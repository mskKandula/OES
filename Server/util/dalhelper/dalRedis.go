package dalhelper

import (
	redis "github.com/go-redis/redis/v8"
)

var redisConnection *redis.Client

func GetRedisConnection(redisDSN string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisDSN)
	if err != nil {
		connectionError = err
	}

	redisConnection = redis.NewClient(opt)

	return redisConnection, connectionError
}
