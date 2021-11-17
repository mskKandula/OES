package websock

import (
	"fmt"
	"log"
)

// var Conns = make(map[*websocket.Conn]bool)

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

			fmt.Println("Sending message to all clients in Pool")

			for _, client := range pool.Clients[message.Id] {

				if client.Details.Role == "Student" {

					if err := client.Conn.WriteJSON(message); err != nil {

						log.Println(err)
						return
					}
				}

			}
		}
	}
}
