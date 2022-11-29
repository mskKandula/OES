package config

import (
	"encoding/json"
	"os"
)

type database struct {
	MySQLDSN    string `json:"MySqlDSN"`
	RedisDSN    string `json:"RedisDSN"`
	RabbitMQDSN string `json:"RabbitMQDSN"`
}

type configuration struct {
	database `json:"database,omitempty"`
}

var DatabaseConfig database

func Setup(fileName string) error {

	var config configuration

	configFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return err
	}

	DatabaseConfig = config.database

	return nil
}
