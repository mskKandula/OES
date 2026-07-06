package dalhelper

import (
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// GetRabbitMQConnection dials RabbitMQ with an explicit TCP dial timeout and
// heartbeat interval, then returns the underlying connection.
// The caller is responsible for opening per-goroutine channels from this
// connection and for closing the connection when done.
func GetRabbitMQConnection(dsn string) (*amqp.Connection, error) {
	return amqp.DialConfig(dsn, amqp.Config{
		// Heartbeat controls how often each side sends a heartbeat frame.
		// If the broker doesn't receive one within 2× this interval it closes
		// the connection, enabling fast detection of dead TCP connections.
		Heartbeat: 10 * time.Second,
		// Dial replaces the default dialer with one that enforces an explicit
		// TCP connection timeout, preventing indefinite hangs at startup.
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 10*time.Second)
		},
	})
}
