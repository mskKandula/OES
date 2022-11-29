package repository

import (
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RepositoryConfig struct {
	MySQLDB  *sql.DB
	Redis    *redis.Client
	RabbitMQ *amqp.Channel
	Queue    amqp.Queue
}
