package websock

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type IntegerArray struct {
	Intarr []int `json : intarr`
}

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
