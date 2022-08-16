package websock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mskKandula/config"
)

// var Conns = make(map[*websocket.Conn]bool)

const PubSubGeneralChannel = "general"

var ctx = context.Background()

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string][]*Client
	Broadcast  chan Message
}

var poolInit *Pool

func init() {
	poolInit = &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string][]*Client),
		Broadcast:  make(chan Message),
	}
}

func NewPool() *Pool {
	return poolInit
}

func (pool *Pool) Start() {

	go pool.listenPubSubChannel()

	for {
		select {

		case client := <-pool.Register:

			pool.Clients[client.Id] = append(pool.Clients[client.Id], client)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients[client.Id]))
			// for _, client := range pool.Clients[client.Id] {

			// 	client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			// }

		case client := <-pool.Unregister:

			delete(pool.Clients, client.Id)

			fmt.Println("Size of Connection Pool : ", len(pool.Clients[client.Id]))

			// for  _,client := range pool.Clients[client.Id] {
			//     client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			// }

		case message := <-pool.Broadcast:
			// Publish the message on "general" channel
			pool.publishMessage(message)
		}
	}
}

// Redis Publish message functionality
func (pool *Pool) publishMessage(msg Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	err = config.Redis.Publish(ctx, PubSubGeneralChannel, payload).Err()

	if err != nil {
		log.Println(err)
		return
	}
}

// Redis Subscribe & Listen on channel("general") functionality
func (pool *Pool) listenPubSubChannel() {

	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)

	msg := Message{}

	ch := pubsub.Channel()

	for data := range ch {

		fmt.Println("Sending message to all clients in Pool")

		if err := json.Unmarshal([]byte(data.Payload), &msg); err != nil {
			log.Println(err)
			return
		}

		for _, client := range pool.Clients[msg.Id] {

			if client.Details.Role == "Student" {

				if err := client.Conn.WriteJSON(msg); err != nil {

					log.Println(err)
					return
				}
			}

		}

	}
}
