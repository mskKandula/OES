package websock

import (
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// const (
// Max wait time when writing message to peer
// writeWait = 10 * time.Second

// Max time till next pong from peer
// pongWait = 60 * time.Second

// Send ping interval, must be less then pong wait time
// pingPeriod = (pongWait * 9) / 10

// Maximum message size allowed from peer.
// maxMessageSize = 100000
// )

// type IntegerArray struct {
// 	Intarr []int `json : intarr`
// }

type Details struct {
	Role string `json:"role"`
	Id   string `json:"id"`
}

type Client struct {
	Conn net.Conn
	Pool *Pool
	*Details
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
	Id   string `json:"id"`
}

func Read(clients <-chan *Client) {
	// defer func() {
	// 	c.Pool.Unregister <- c
	// 	c.Conn.Close()
	// }()

	// c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// m := Message{}

	// err := c.Conn.ReadJSON(&m)
	// if err != nil {
	// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
	// 		log.Printf("unexpected close error: %v", err)
	// 	}
	// 	log.Println(err)
	// 	break
	// }

	for c := range clients {

		msg, op, err := wsutil.ReadClientData(c.Conn)
		if err != nil {
			log.Println("unknown message", "error", err.Error())
			return
		}

		switch op {
		case ws.OpPing:
			if err := wsutil.WriteServerMessage(c.Conn, ws.OpPong, msg); err != nil {
				log.Println("failed to write pong message ", "error", err.Error())
				return
			}
		case ws.OpText:
			log.Println("receving message from client", "message", msg)
			c.Pool.Broadcast <- msg
		case ws.OpClose:
			log.Println("client has been disconnected")
			c.Pool.Unregister <- c
			c.Conn.Close()
		}

		// if err = json.Unmarshal(byteData, &m); err != nil {
		// 	log.Println(err)
		// 	return
		// }

		// if m.Type == 2 {
		// 	chart(c)
		// } else {
		// 	c.Pool.Broadcast <- m
		// }
		// fmt.Printf("Message Received: %+v\n", m)
	}
}

// func chart(c *Client) {

// 	for i := 0; i < 5; i++ {

// 		rand.Seed(time.Now().UnixNano())

// 		i := &IntegerArray{

// 			Intarr: rand.Perm(12),
// 		}

// 		byteData, err := json.Marshal(i)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		if err := wsutil.WriteClientBinary(c.Conn, byteData); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		time.Sleep(3 * time.Second)
// 	}
// }
