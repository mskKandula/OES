package ds

import (
	"database/sql"
	"fmt"

	"github.com/mskKandula/oes/api/config"
	"github.com/mskKandula/oes/util/dalhelper"

	redis "github.com/go-redis/redis/v8"
)

type DataSources struct {
	MySQLDB *sql.DB
	Redis   *redis.Client
}

//InitDS initializes  all required data sources.
//Returns MySQL and MongoDB  client instances as DataSource
func InitDS() (*DataSources, error) {
	mySQLDB, err := dalhelper.GetMySQLConnection(config.DatabaseConfig.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening mysqldb : %w", err)
	}

	redis, err := dalhelper.GetRedisConnection(config.DatabaseConfig.RedisDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening redis : %w", err)
	}

	return &DataSources{
		MySQLDB: mySQLDB,
		Redis:   redis,
	}, nil
}
