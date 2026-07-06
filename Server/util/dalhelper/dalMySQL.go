package dalhelper

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Hold a single global connection (pooling provided by sql driver)
var sqlConnection *sql.DB
var connectionError error
var sqlOnce sync.Once

// GetMySQLConnection creates a singleton MySQL connection pool.
// Returns the connection if established, otherwise returns an error.
func GetMySQLConnection(mySqlDSN string) (*sql.DB, error) {
	// sqlOnce ensures only one connection pool is ever created.
	sqlOnce.Do(func() {
		// sql.Open is lazy — it validates the DSN but does NOT dial the server.
		// Errors here are DSN parse failures only, not connectivity failures.
		connection, err := sql.Open("mysql", mySqlDSN)
		if err != nil {
			connectionError = err
			return // guard: do not configure a nil connection
		}

		// Optimized connection pool settings for better performance.
		// Increased idle connections to reduce connection overhead.
		connection.SetMaxIdleConns(50)
		// Increased max open connections for better concurrency.
		connection.SetMaxOpenConns(200)
		// Reduced idle timeout to prevent stale connections.
		connection.SetConnMaxIdleTime(10 * time.Minute)
		// Connection max lifetime to ensure fresh connections.
		connection.SetConnMaxLifetime(1 * time.Hour)

		// Ping verifies the server is reachable and the DSN is valid.
		// Without this, the app would start successfully even when MySQL is
		// unreachable, surfacing the error only on the first query.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := connection.PingContext(ctx); err != nil {
			connectionError = fmt.Errorf("mysql ping failed: %w", err)
			return
		}

		sqlConnection = connection
	})
	return sqlConnection, connectionError
}
