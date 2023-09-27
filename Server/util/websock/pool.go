package websock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
	ds "github.com/mskKandula/oes/dataSources"
)

// var Conns = make(map[*websocket.Conn]bool)

const PubSubGeneralChannel = "general"

var (
	ctx            = context.Background()
	poolInit       *Pool
	ClientConnChan chan *Client
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string][]*Client
	Broadcast  chan []byte
}

func init() {
	poolInit = &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string][]*Client),
		Broadcast:  make(chan []byte),
	}
	ClientConnChan = make(chan *Client, 200)
}

func NewPool() *Pool {
	return poolInit
}

func (pool *Pool) Start(ds *ds.DataSources) {

	go pool.listenPubSubChannel(ds)
	poller, err := netpoll.New(nil)
	if err != nil {
		log.Println(err)
	}

	for {
		select {

		case newClient := <-pool.Register:

			pool.Clients[newClient.Id] = append(pool.Clients[newClient.Id], newClient)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients[newClient.Id]))

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

			// for _, client := range pool.Clients[client.Id] {

			byteData, err := json.Marshal(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": true}, Id: newClient.Id})
			if err != nil {
				log.Println(err)
				return
			}

			for _, client := range pool.Clients[newClient.Id] {

				if client.Role != newClient.Role || client.UserId != newClient.UserId {

					if err := wsutil.WriteServerText(client.Conn, byteData); err != nil {
						log.Println(err)
						return
					}

					// client.Conn.WriteJSON(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": true}, Id: newClient.Id})
				}

			}

		case newClient := <-pool.Unregister:

			var toDelete int
			byteData, err := json.Marshal(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": false}, Id: newClient.Id})
			if err != nil {
				log.Println(err)
				return
			}

			for i, client := range pool.Clients[newClient.Id] {

				if client.UserId == newClient.UserId {
					toDelete = i
				}

				if client.Role != newClient.Role || client.UserId != newClient.UserId {
					if err := wsutil.WriteServerText(client.Conn, byteData); err != nil {
						log.Println(err)
						return
					}

					// client.Conn.WriteJSON(Message{Type: 5, Body: map[string]interface{}{"user": newClient.UserId, "add": false}, Id: newClient.Id})
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
func (pool *Pool) publishMessage(payload []byte, ds *ds.DataSources) {
	// payload, err := json.Marshal(msg)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	err := ds.Redis.Publish(ctx, PubSubGeneralChannel, payload).Err()

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
				if err := wsutil.WriteServerText(client.Conn, []byte(data.Payload)); err != nil {
					log.Println(err)
					return
				}

			} else if msg.Type == 4 && client.UserId == msg.To {
				fmt.Println("Sending message to specific client in Pool")
				if err := wsutil.WriteServerText(client.Conn, []byte(data.Payload)); err != nil {
					log.Println(err)
					return
				}
			}

		}

	}
}
