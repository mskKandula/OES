package websock

import (
	"encoding/json"
	"log"
	"net"
	"sync/atomic"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Details struct {
	Role   string `json:"role"`
	Id     string `json:"id"`
	UserId int    `json:"userId"`
}

type Client struct {
	Conn net.Conn
	Pool *Pool
	*Details
	IsAlive atomic.Bool
}

type JoinLeaveBody struct {
	User int  `json:"user"`
	Add  bool `json:"add"`
}
type Message struct {
	Type int             `json:"type"`
	Body json.RawMessage `json:"body"`
	Id   string          `json:"id"`
	To   int             `json:"to,omitempty"`
}

func Read(clients <-chan *Client) {

	for c := range clients {

		msg, op, err := wsutil.ReadClientData(c.Conn)
		if err != nil {
			log.Println("unknown message", "error", err.Error())
			continue
		}

		switch op {

		case ws.OpPong:
			c.IsAlive.Store(true)

		case ws.OpText:
			c.IsAlive.Store(true)
			log.Println("receving message from client", "message", msg)
			c.Pool.Broadcast <- msg
		case ws.OpClose:
			log.Println("client has been disconnected")
			c.Pool.Unregister <- c
		}
	}
}
