package websock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gobwas/ws"
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

		case client := <-pool.Register:

			pool.Clients[client.Id] = append(pool.Clients[client.Id], client)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients[client.Id]))

			// Get netpoll descriptor with EventRead|EventEdgeTriggered.
			desc := netpoll.Must(netpoll.HandleRead(client.Conn))

			// Make conn to be observed by netpoll instance.
			poller.Start(desc, func(ev netpoll.Event) {
				if ev&netpoll.EventReadHup != 0 {
					poller.Stop(desc)
					client.Conn.Close()
					return
				}
				ClientConnChan <- client
			})

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

				if err := wsutil.WriteServerMessage(client.Conn, ws.OpText, []byte(data.Payload)); err != nil {

					log.Println(err)
					return
				}
			}

		}

	}
}
