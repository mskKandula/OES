package dalhelper

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// GetRedisConnection parses the DSN, configures a connection pool,
// pings the server to verify connectivity, and returns the client.
func GetRedisConnection(dsn string) (*redis.Client, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	opt.PoolSize = 10
	opt.MinIdleConns = 2
	opt.MaxConnAge = 30 * time.Minute
	opt.PoolTimeout = 4 * time.Second
	opt.IdleTimeout = 5 * time.Minute
	opt.IdleCheckFrequency = 1 * time.Minute

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
