package websock

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
)

// var Conns = make(map[*websocket.Conn]bool)

const (
	PubSubGeneralChannel        = "general"
	HeartBeatInterval           = 10 * time.Second
	ShardCount           uint32 = 32 // Number of shards to divide the map into
)

var (
	ctx            = context.Background()
	poolInit       *Pool
	ClientConnChan chan *Client
)

// PoolMapShard holds its own Mutex and Map to reduce lock contention
type PoolMapShard struct {
	sync.RWMutex
	Clients map[string]map[int]*Client
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Shards     []*PoolMapShard
	Broadcast  chan []byte
}

func getShardIndex(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32() % ShardCount
}

func (pool *Pool) getShard(id string) *PoolMapShard {
	return pool.Shards[getShardIndex(id)]
}

func init() {
	poolInit = &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Shards:     make([]*PoolMapShard, ShardCount),
		Broadcast:  make(chan []byte, 256),
	}

	for i := 0; i < int(ShardCount); i++ {
		poolInit.Shards[i] = &PoolMapShard{
			Clients: make(map[string]map[int]*Client),
		}
	}

	ClientConnChan = make(chan *Client, 200)
}

func NewPool() *Pool {
	return poolInit
}

func (pool *Pool) Start(redis *redis.Client) {

	// Tick faster, but process only one shard per tick to spread the load
	ticker := time.NewTicker(HeartBeatInterval / time.Duration(ShardCount))
	var currentShardIndex uint32 = 0

	go pool.listenPubSubChannel(redis)
	poller, err := netpoll.New(nil)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		close(ClientConnChan)
	}()

	for {
		select {

		case newClient := <-pool.Register:

			shard := pool.getShard(newClient.Id)
			shard.Lock()
			if shard.Clients[newClient.Id] == nil {
				shard.Clients[newClient.Id] = make(map[int]*Client)
			}

			shard.Clients[newClient.Id][newClient.UserId] = newClient
			poolSize := len(shard.Clients[newClient.Id])
			shard.Unlock()

			fmt.Println("Size of Connection Pool: ", poolSize)

			// Get netpoll descriptor with EventRead|EventEdgeTriggered.
			desc := netpoll.Must(netpoll.HandleRead(newClient.Conn))

			// Make conn to be observed by netpoll instance.
			poller.Start(desc, func(ev netpoll.Event) {
				if ev&netpoll.EventReadHup != 0 {
					poller.Stop(desc)
					newClient.Conn.Close()
					return
				}
				ClientConnChan <- newClient
			})

			rawBody, _ := json.Marshal(JoinLeaveBody{User: newClient.UserId, Add: true})
			byteData, err := json.Marshal(Message{Type: 5, Body: json.RawMessage(rawBody), Id: newClient.Id})
			if err != nil {
				log.Println(err)
				return
			}

			pool.Broadcast <- byteData

		case newClient := <-pool.Unregister:

			rawBody, _ := json.Marshal(JoinLeaveBody{User: newClient.UserId, Add: false})
			byteData, err := json.Marshal(Message{Type: 5, Body: json.RawMessage(rawBody), Id: newClient.Id})
			if err != nil {
				log.Println(err)
				return
			}

			pool.Broadcast <- byteData

			shard := pool.getShard(newClient.Id)
			shard.Lock()
			delete(shard.Clients[newClient.Id], newClient.UserId)

			poolSize := 0
			if len(shard.Clients[newClient.Id]) == 0 {
				delete(shard.Clients, newClient.Id)
			} else {
				poolSize = len(shard.Clients[newClient.Id])
			}
			shard.Unlock()

			// Close the connection
			newClient.Conn.Close()

			fmt.Println("Size of Connection Pool: ", poolSize)

		case message := <-pool.Broadcast:
			// Publish the message on "general" channel
			pool.publishMessage(message, redis)

		case <-ticker.C:
			// Process clients incrementally to reduce memory pressure
			// Process one shard per tick to perfectly spread the load over the HeartBeatInterval
			shard := pool.Shards[currentShardIndex]
			shard.RLock()
			for _, clientsMap := range shard.Clients {
				for _, client := range clientsMap {
					// Check and update IsAlive status
					if !client.IsAlive.Load() {
						shard.RUnlock()
						client.Pool.Unregister <- client
						shard.RLock()
						continue
					}

					client.IsAlive.Store(false)

					if err := ws.WriteFrame(client.Conn, ws.NewPingFrame(nil)); err != nil {
						log.Println(err)
						continue
					}
				}
			}
			shard.RUnlock()

			// Move to the next shard for the next tick
			currentShardIndex = (currentShardIndex + 1) % ShardCount
		}
	}
}

// Redis Publish message functionality
func (pool *Pool) publishMessage(payload []byte, redis *redis.Client) {

	err := redis.Publish(ctx, PubSubGeneralChannel, payload).Err()

	if err != nil {
		log.Println(err)
		return
	}
}

// Redis Subscribe & Listen on channel("general") functionality
func (pool *Pool) listenPubSubChannel(redis *redis.Client) {

	pubsub := redis.Subscribe(ctx, PubSubGeneralChannel)

	ch := pubsub.Channel()

	for data := range ch {

		msg := Message{}
		payloadData := []byte(data.Payload)

		if err := json.Unmarshal(payloadData, &msg); err != nil {
			log.Println(err)
			return
		}

		shard := pool.getShard(msg.Id)
		shard.RLock()

		switch msg.Type {
		case 1, 3:
			for _, client := range shard.Clients[msg.Id] {
				fmt.Println("Sending message to all students in Pool")
				if client.Details.Role == "Student" {
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println(err)
						continue
					}
				}
			}

		case 4:
			if client, ok := shard.Clients[msg.Id][msg.To]; ok {
				fmt.Println("Sending message to specific client in Pool")
				if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
					log.Println(err)
					continue
				}
			}

		case 5:
			var body JoinLeaveBody
			json.Unmarshal(msg.Body, &body)
			joiningUserId := body.User
			for userId, client := range shard.Clients[msg.Id] {
				if userId != joiningUserId {
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}

		shard.RUnlock()

	}
}
