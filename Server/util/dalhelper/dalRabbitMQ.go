package dalhelper

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements model.Publisher backed by an AMQP channel.
type RabbitMQPublisher struct {
	ch *amqp.Channel
}

// NewRabbitMQPublisher wraps an *amqp.Channel in the Publisher interface.
func NewRabbitMQPublisher(ch *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{ch: ch}
}

func GetRabbitMQConnectionWithQueues(rabbitMQDSN string) (*amqp.Channel, error) {
	conn, connError := amqp.Dial(rabbitMQDSN)
	if connError != nil {
		return nil, connError
	}

	ch, err := conn.Channel()
	if err != nil {
		return ch, err
	}
	_, err = DeclareEncodeQueue(ch)
	if err != nil {
		return ch, err
	}
	_, err = DeclareEmailQueue(ch)
	if err != nil {
		return ch, err
	}

	return ch, nil
}

// DeclareEncodeQueue declares the durable "encode" queue on the given channel.
// Using durable:true so pending encode jobs survive a RabbitMQ broker restart.
func DeclareEncodeQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"encode", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
}

// DeclareEmailQueue declares the durable "email" queue on the given channel.
// Using durable:true so pending email jobs survive a RabbitMQ broker restart.
func DeclareEmailQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"email", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
}

// PublishMessage publishes a JSON body to the named queue as a persistent message.
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
