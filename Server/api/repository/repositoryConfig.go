package repository

import (
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type RepositoryConfig struct {
	MySQLDB *sql.DB
	Redis   *redis.Client
}
