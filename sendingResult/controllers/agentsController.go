package controllers

import (
	"encoding/json"
	"log"
)

type AgentsAuthData struct {
	Url      string
	Login    string
	Password string
}

type AgentsData struct {
	Url            string
	AgentsAuthData AgentsAuthData
	JsonToSend     string `json:"data"`
}

//AgentsController is a root structure needed for functions invocation through reflection. This functions handle results of data sending to outer api's
type AgentsController struct{}

func parseData(data *json.RawMessage) (*AgentsData, error) {
	dataObject := &AgentsData{}

	err := json.Unmarshal(*data, dataObject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dataObject, nil
}

func (agents *AgentsController) HandleSendOrderSuccess(data *json.RawMessage) error {
	agentData, err := parseData(data)

	if err != nil {
		return err
	}

	log.Println(agentData) //work imitation

	return nil
}

func (agents *AgentsController) HandleSendOrderFailure(data *json.RawMessage) error {
	agentData, err := parseData(data)

	if err != nil {
		return err
	}

	log.Println(agentData) //work imitation

	return nil
}

func (agents *AgentsController) HandleSendCancelSuccess(data *json.RawMessage) error {
	agentData, err := parseData(data)

	if err != nil {
		return err
	}

	log.Println(agentData) //work imitation

	return nil
}

func (agents *AgentsController) HandleSendCancelFailure(data *json.RawMessage) error {
	agentData, err := parseData(data)

	if err != nil {
		return err
	}

	log.Println(agentData) //work imitation

	return nil
}
