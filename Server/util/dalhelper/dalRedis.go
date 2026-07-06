package dalhelper

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

//sudo service redis-server start

var redisConnection *redis.Client

func GetRedisConnection(redisDSN string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisDSN)
	if err != nil {
		connectionError = err
		return nil, connectionError
	}

	// Optimize Redis connection pool settings
	opt.PoolSize = 50                        // Increased pool size for better concurrency
	opt.MinIdleConns = 10                    // Keep minimum idle connections ready
	opt.MaxConnAge = 30 * time.Minute        // Recycle connections periodically
	opt.PoolTimeout = 4 * time.Second        // Timeout waiting for connection from pool
	opt.IdleTimeout = 5 * time.Minute        // Close idle connections after timeout
	opt.IdleCheckFrequency = 1 * time.Minute // Check for idle connections frequency

	redisConnection = redis.NewClient(opt)

	// Fail fast at startup: verify Redis is reachable and accepting commands.
	// redis.NewClient is lazy — it does not dial until the first command is issued,
	// so without this Ping the app would start successfully even when Redis is
	// unreachable or misconfigured, only surfacing the error on the first request.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisConnection.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return redisConnection, nil
}
