package appconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	RabbitMq   RabbitMq `json:"rabbitMq"`
	PipelineID int      `json:"pipelineId"`
}

type RabbitMq struct {
	ConnectionString string `json:"connectionString"`
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
