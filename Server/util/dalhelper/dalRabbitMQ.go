package dalhelper

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ch        *amqp.Channel
	connError error
	q         amqp.Queue
)

func GetRabbitMQConnection(rabbitMQDSN string) (*amqp.Channel, amqp.Queue, error) {
	conn, connError := amqp.Dial(rabbitMQDSN)
	if connError != nil {
		return ch, q, connError
	}

	ch, connError = conn.Channel()
	if connError != nil {
		return ch, q, connError
	}

	q, connError = ch.QueueDeclare(
		"encode", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if connError != nil {
		return ch, q, connError
	}
	return ch, q, connError
}
