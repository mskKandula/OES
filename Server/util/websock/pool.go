package websock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	ds "github.com/mskKandula/oes/dataSources"
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

func (pool *Pool) Start(ds *ds.DataSources) {

	go pool.listenPubSubChannel(ds)

	for {
		select {

		case newClient := <-pool.Register:

			pool.Clients[newClient.Id] = append(pool.Clients[newClient.Id], newClient)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients[newClient.Id]))

			for _, client := range pool.Clients[newClient.Id] {

				if client.Role != newClient.Role || client.UserId != newClient.UserId {
					client.Conn.WriteJSON(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": true}, Id: newClient.Id})
				}

			}

		case newClient := <-pool.Unregister:

			var toDelete int

			for i, client := range pool.Clients[newClient.Id] {

				if client.UserId == newClient.UserId {
					toDelete = i
				}

				if client.Role != newClient.Role || client.UserId != newClient.UserId {
					client.Conn.WriteJSON(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": false}, Id: newClient.Id})
				}

			}
			pool.Clients[newClient.Id] = append(pool.Clients[newClient.Id][:toDelete], pool.Clients[newClient.Id][toDelete+1:]...)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients[newClient.Id]))

		case message := <-pool.Broadcast:
			// Publish the message on "general" channel
			pool.publishMessage(message, ds)
		}
	}
}

// Redis Publish message functionality
func (pool *Pool) publishMessage(msg Message, ds *ds.DataSources) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	err = ds.Redis.Publish(ctx, PubSubGeneralChannel, payload).Err()

	if err != nil {
		log.Println(err)
		return
	}
}

// Redis Subscribe & Listen on channel("general") functionality
func (pool *Pool) listenPubSubChannel(ds *ds.DataSources) {

	pubsub := ds.Redis.Subscribe(ctx, PubSubGeneralChannel)

	ch := pubsub.Channel()

	for data := range ch {

		msg := Message{}

		if err := json.Unmarshal([]byte(data.Payload), &msg); err != nil {
			log.Println(err)
			return
		}

		for _, client := range pool.Clients[msg.Id] {
			if msg.Type == 1 && client.Details.Role == "Student" {
				fmt.Println("Sending message to all students in Pool")
				if err := client.Conn.WriteJSON(msg); err != nil {

					log.Println(err)
					return
				}

			} else if msg.Type == 4 && client.UserId == msg.To {
				fmt.Println("Sending message to specific client in Pool")
				if err := client.Conn.WriteJSON(msg); err != nil {

					log.Println(err)
					return
				}
			}

		}

	}
}
