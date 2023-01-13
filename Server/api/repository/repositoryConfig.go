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

// type Transaction interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	Prepare(query string) (*sql.Stmt, error)
// 	Query(query string, args ...interface{}) (*sql.Rows, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// }

type TxFn func(*sql.Tx) error

//Not an exported one.
func withTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
