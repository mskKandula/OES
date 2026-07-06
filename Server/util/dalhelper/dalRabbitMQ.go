package dalhelper

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// queues lists every durable queue the server needs to declare on startup.
// All are durable (survive broker restart) and non-exclusive (shared across connections).
var queues = []string{
	"encode",
	"email",
	"code.execute.python",
	"code.execute.go",
	"code.execute.nodejs",
}

// RabbitMQPublisher implements model.Publisher backed by an AMQP channel.
type RabbitMQPublisher struct {
	ch *amqp.Channel
}

// NewRabbitMQPublisher wraps an *amqp.Channel in the Publisher interface.
func NewRabbitMQPublisher(ch *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{ch: ch}
}

// GetRabbitMQConnectionWithQueues dials RabbitMQ and idempotently declares all
// required durable queues before returning the channel.
func GetRabbitMQConnectionWithQueues(rabbitMQDSN string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitMQDSN)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, name := range queues {
		if _, err = declareQueue(ch, name); err != nil {
			return nil, fmt.Errorf("declare queue %q: %w", name, err)
		}
	}

	return ch, nil
}

// declareQueue declares a single durable, non-exclusive, non-auto-delete queue.
// It is idempotent — safe to call even if the queue already exists.
func declareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// PublishMessageWithContext publishes a JSON body to the named queue as a persistent message.
func (p *RabbitMQPublisher) PublishMessageWithContext(ctx context.Context, queueName string, body []byte) error {
	return p.ch.PublishWithContext(
		ctx,
		"",        // default exchange
		queueName, // routing key == queue name
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}
