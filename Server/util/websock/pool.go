package websock

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
)

const (
	PubSubGeneralChannel        = "general"
	HeartBeatInterval           = 10 * time.Second
	ShardCount           uint32 = 32 // Number of shards to divide the map into

	// WriteDeadline is the maximum time allowed for a single WebSocket write.
	WriteDeadline = 100 * time.Millisecond

	// shardWorkerBufSize is the per-shard task channel buffer.
	shardWorkerBufSize = 256
)

var (
	ctx            = context.Background()
	poolInit       *Pool
	ClientConnChan chan *Client

	// shardWorkers is a fixed-size array of buffered channels — one per shard.
	shardWorkers [ShardCount]chan shardTask
)

// PoolMapShard holds its own Mutex and Map to reduce lock contention
type PoolMapShard struct {
	sync.RWMutex
	Clients map[string]map[int]*Client
}

// shardTask is the unit of work dispatched to a per-shard worker goroutine.
type shardTask struct {
	msg         Message
	payloadData []byte
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
	for i := range shardWorkers {
		shardWorkers[i] = make(chan shardTask, shardWorkerBufSize)
	}

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

	// Start one worker goroutine per shard. Each worker owns the write loop
	// for its shard, so all 32 shards can flush WebSocket writes concurrently.
	for i := uint32(0); i < ShardCount; i++ {
		go pool.runShardWorker(i, shardWorkers[i])
	}

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

			log.Println("Size of Connection Pool: ", poolSize)

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

			log.Println("Size of Connection Pool: ", poolSize)

		case message := <-pool.Broadcast:
			// Publish the message on "general" channel
			pool.publishMessage(message, redis)

		case <-ticker.C:
			// Dispatch the ping work for this shard to a separate goroutine
			go pool.pingShardClients(currentShardIndex)

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
			continue
		}

		shardIdx := getShardIndex(msg.Id)

		// Non-blocking send: if the worker's buffer is full the message is
		// dropped for that shard rather than stalling the listener goroutine
		// and backing up every other shard.
		select {
		case shardWorkers[shardIdx] <- shardTask{msg: msg, payloadData: payloadData}:
		default:
			log.Printf("shard %d worker channel full, dropping message type %d for id %s", shardIdx, msg.Type, msg.Id)
		}
	}
}

// pingShardClients sends a WebSocket ping frame to every client in the given
// shard and checks their IsAlive status.One goroutine is spawned per tick per shard 
func (pool *Pool) pingShardClients(shardIndex uint32) {
	shard := pool.Shards[shardIndex]
	shard.RLock()

	// Collect dead clients first so we don't call Unregister while holding.
	var dead []*Client

	for _, clientsMap := range shard.Clients {
		for _, client := range clientsMap {
			if !client.IsAlive.Load() {
				dead = append(dead, client)
				continue
			}

			client.IsAlive.Store(false)

			client.Conn.SetWriteDeadline(time.Now().Add(WriteDeadline))
			if err := ws.WriteFrame(client.Conn, ws.NewPingFrame(nil)); err != nil {
				log.Println("ping error:", err)
			}
		}
	}

	shard.RUnlock()

	// Enqueue dead clients for cleanup outside the read lock.
	for _, client := range dead {
		client.Pool.Unregister <- client
	}
}

// runShardWorker is the dedicated write goroutine for a single shard.
func (pool *Pool) runShardWorker(idx uint32, tasks <-chan shardTask) {
	shard := pool.Shards[idx]

	for task := range tasks {
		msg := task.msg
		payloadData := task.payloadData

		shard.RLock()

		switch msg.Type {
		case 1, 3:
			for _, client := range shard.Clients[msg.Id] {
				if client.Details.Role == "Student" {
					log.Println("Sending message to all students in Pool")
					client.Conn.SetWriteDeadline(time.Now().Add(WriteDeadline))
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println("shard worker write error:", err)
						continue
					}
				}
			}

		case 4:
			if client, ok := shard.Clients[msg.Id][msg.To]; ok {
				log.Println("Sending message to specific client in Pool")
				client.Conn.SetWriteDeadline(time.Now().Add(WriteDeadline))
				if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
					log.Println("shard worker write error:", err)
				}
			}

		case 5:
			var body JoinLeaveBody
			json.Unmarshal(msg.Body, &body)
			joiningUserId := body.User
			for userId, client := range shard.Clients[msg.Id] {
				if userId != joiningUserId {
					client.Conn.SetWriteDeadline(time.Now().Add(WriteDeadline))
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println("shard worker write error:", err)
						continue
					}
				}
			}
		}

		shard.RUnlock()
	}
}
