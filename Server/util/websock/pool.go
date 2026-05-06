package websock

import (
	"context"
	"encoding/json"
	"fmt"
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
	PubSubGeneralChannel = "general"
	HeartBeatInterval    = 10 * time.Second
)

var (
	ctx            = context.Background()
	poolInit       *Pool
	ClientConnChan chan *Client
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]map[int]*Client
	Broadcast  chan []byte
	mu         sync.RWMutex // Protects Clients map
}

func init() {
	poolInit = &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]map[int]*Client),
		Broadcast:  make(chan []byte, 256),
	}
	ClientConnChan = make(chan *Client, 200)
}

func NewPool() *Pool {
	return poolInit
}

func (pool *Pool) Start(redis *redis.Client) {

	ticker := time.NewTicker(HeartBeatInterval)

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

			// Initialize nested map if it doesn't exist
			pool.mu.Lock()
			if pool.Clients[newClient.Id] == nil {
				pool.Clients[newClient.Id] = make(map[int]*Client)
			}

			// Add client using UserId as key in the nested map
			pool.Clients[newClient.Id][newClient.UserId] = newClient
			poolSize := len(pool.Clients[newClient.Id])
			pool.mu.Unlock()

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

			// Delete client from nested map using delete() function
			pool.mu.Lock()
			delete(pool.Clients[newClient.Id], newClient.UserId)

			// Clean up empty nested maps
			poolSize := 0
			if len(pool.Clients[newClient.Id]) == 0 {
				delete(pool.Clients, newClient.Id)
			} else {
				poolSize = len(pool.Clients[newClient.Id])
			}
			pool.mu.Unlock()

			// Close the connection
			newClient.Conn.Close()

			fmt.Println("Size of Connection Pool: ", poolSize)

		case message := <-pool.Broadcast:
			// Publish the message on "general" channel
			pool.publishMessage(message, redis)

		case <-ticker.C:
			// Process clients incrementally to reduce memory pressure
			pool.mu.RLock()
			for _, clientsMap := range pool.Clients {
				for _, client := range clientsMap {
					// Check and update IsAlive status
					if !client.IsAlive.Load() {
						// Unlock before sending to channel to avoid blocking
						pool.mu.RUnlock()
						client.Pool.Unregister <- client
						pool.mu.RLock()
						continue
					}

					client.IsAlive.Store(false)

					// Send heartbeat - do this while holding read lock is safe
					// as we're not modifying the map structure
					// Use:
					if err := ws.WriteFrame(client.Conn, ws.NewPingFrame(nil)); err != nil {
						log.Println(err)
						continue
					}
				}
			}
			pool.mu.RUnlock()
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

		pool.mu.RLock()

		switch msg.Type {
		case 1, 3:
			for _, client := range pool.Clients[msg.Id] {
				fmt.Println("Sending message to all students in Pool")
				if client.Details.Role == "Student" {
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println(err)
						continue
					}
				}
			}

		case 4:
			if client, ok := pool.Clients[msg.Id][msg.To]; ok {
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
			for userId, client := range pool.Clients[msg.Id] {
				if userId != joiningUserId {
					if err := wsutil.WriteServerText(client.Conn, payloadData); err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}

		pool.mu.RUnlock()

	}
}
