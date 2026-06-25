package repository

import (
	"context"
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type RepositoryConfig struct {
	MySQLDB *sql.DB
	Redis   *redis.Client
}

// type Transaction interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	Prepare(query string) (*sql.Stmt, error)
// 	Query(query string, args ...interface{}) (*sql.Rows, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// }

type TxFn func(*sql.Tx) error

// Not an exported one.
// withTransactionContext wraps database operations in a transaction with context support for better control
func withTransactionContext(ctx context.Context, db *sql.DB, fn TxFn) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
