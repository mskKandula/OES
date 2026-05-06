package ds

import (
	"database/sql"
	"fmt"

	"github.com/mskKandula/oes/api/config"
	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/util/dalhelper"

	redis "github.com/go-redis/redis/v8"
)

type DataSources struct {
	MySQLDB   *sql.DB
	Redis     *redis.Client
	Publisher model.Publisher
}

// InitDS initializes all required data sources.
// Returns MySQL, Redis, and a RabbitMQ-backed Publisher as DataSources.
func InitDS() (*DataSources, error) {
	mySQLDB, err := dalhelper.GetMySQLConnection(config.DatabaseConfig.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening mysqldb : %w", err)
	}

	redis, err := dalhelper.GetRedisConnection(config.DatabaseConfig.RedisDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening redis : %w", err)
	}

	rabbitMQCh, err := dalhelper.GetRabbitMQConnectionWithQueues(config.DatabaseConfig.RabbitMQDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening rabbitMQ : %w", err)
	}

	publisher := dalhelper.NewRabbitMQPublisher(rabbitMQCh)

	return &DataSources{
		MySQLDB:   mySQLDB,
		Redis:     redis,
		Publisher: publisher,
	}, nil
}
