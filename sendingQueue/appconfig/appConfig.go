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

type AppConfig struct {
	MongoConfig MongoConfig `json:"mongo"`
	RabbitMq    RabbitMq    `json:"rabbitMq"`
	PipelineID  int         `json:"pipelineId"`
}

type RabbitMq struct {
	ConnectionString string `json:"connectionString"`
}

func (config *AppConfig) configureConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Config file reading error", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("Config json deserialization error", err)
	}
}

func NewConfig() AppConfig {
	var appConfig = AppConfig{}
	appConfig.configureConfig()
	return appConfig
}
