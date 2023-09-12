package websock

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	// writeWait = 10 * time.Second

	// Max time till next pong from peer
	// pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	// pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 100000
)

type IntegerArray struct {
	Intarr []int `json : intarr`
}

type Details struct {
	Role   string `json:"role"`
	Id     string `json:"id"`
	UserId int    `json:"userId"`
}

type Client struct {
	Conn *websocket.Conn
	Pool *Pool
	*Details
}

type Message struct {
	Type int                    `json:"type"`
	Body map[string]interface{} `json:"body"`
	Id   string                 `json:"id"`
	To   int                    `json:"to,omitempty"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		m := Message{}

		err := c.Conn.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			log.Println(err)
			break
		}

		if m.Type == 2 {
			chart(c)
		} else {
			c.Pool.Broadcast <- m
		}

		// fmt.Printf("Message Received: %+v\n", m)
	}
}

func chart(c *Client) {

	for i := 0; i < 5; i++ {

		rand.Seed(time.Now().UnixNano())

		i := &IntegerArray{

			Intarr: rand.Perm(12),
		}

		if err := c.Conn.WriteJSON(i); err != nil {
			log.Println(err)
			return
		}

		time.Sleep(3 * time.Second)
	}
	return
}
