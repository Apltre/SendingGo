package appconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type MongoConfig struct {
	ConnectionString string `json:"connectionString"`
	CommandTimeout   int    `json:"commandTimeout"`
}

type APIConfig struct {
	Port int `json:"port"`
}

type AppConfig struct {
	MongoConfig MongoConfig `json:"mongo"`
	APIConfig   APIConfig   `json:"api"`
}

func (config *AppConfig) configureConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Config file reading error: %s\n", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("Config json deserialization error: %s\n", err)
	}
}

func NewConfig() AppConfig {
	var appConfig = AppConfig{}
	appConfig.configureConfig()
	return appConfig
}
