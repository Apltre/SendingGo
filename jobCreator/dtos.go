package main

import "time"

//JobDto is expected data structure received from http request to this service
type JobDto struct {
	Status      int         `json:"status"`
	OperationID int         `json:"operationId"`
	Data        interface{} `json:"data"`
	PipelineID  int         `json:"pipelineId"`
	StartAt     *time.Time  `json:"startAt"`
}
