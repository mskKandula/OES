package websock

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Details struct {
	Role string `json:"role"`
	Id   string `json:"id"`
}

type Client struct {
	Conn *websocket.Conn
	Pool *Pool
	*Details
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
	Id   string `json:"id"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		m := Message{}

		err := c.Conn.ReadJSON(&m)

		if err != nil {
			log.Println(err)
			return
		}

		c.Pool.Broadcast <- m

		fmt.Printf("Message Received: %+v\n", m)
	}
}
