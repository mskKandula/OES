package dalhelper

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Hold a single global connection (pooling provided by sql driver)
var sqlConnection *sql.DB
var connectionError error
var sqlOnce sync.Once

//GetMySQLConnection creates MySQLConnection.
//Returns connection if established else error
func GetMySQLConnection(mySqlDSN string) (*sql.DB, error) {
	//sqlOnce is used to create singleton object
	sqlOnce.Do(func() {
		// create a connection db(e.g. "postgres", "mysql", or "sqlite3")
		connection, err := sql.Open("mysql", mySqlDSN)
		if err != nil {
			connectionError = err
		}
		// Optimized connection pool settings for better performance
		// Increased idle connections to reduce connection overhead
		connection.SetMaxIdleConns(50)
		// Increased max open connections for better concurrency
		connection.SetMaxOpenConns(200)
		// Reduced idle timeout to prevent stale connections
		connection.SetConnMaxIdleTime(10 * time.Minute)
		// Connection max lifetime to ensure fresh connections
		connection.SetConnMaxLifetime(1 * time.Hour)
		sqlConnection = connection
	})
	return sqlConnection, connectionError
}
