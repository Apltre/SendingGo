package models

import (
	"encoding/json"
	"sendingQueue/data/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID `json:"id"`
	OperationID         int                `json:"operationID"`
	Data                *json.RawMessage   `json:"data"`
	Status              int                `json:"status"`
	SendingAttemptIndex int                `json:"sendingAttemptIndex"`
	SendingError        *string            `json:"sendingError"`
	PipelineID          int                `json:"pipelineID"`
	StartAt             time.Time          `json:"startAt"`
	ProcessedDate       *time.Time         `json:"processedDate"`
}

func MongoJobToModel(mongoJob entities.Job) Job {
	var rawDataPtr *json.RawMessage
	if mongoJob.Data != nil {
		var rawData = json.RawMessage(mongoJob.Data.String())
		rawDataPtr = &rawData
	}
	return Job{
		ID:                  mongoJob.ID,
		OperationID:         mongoJob.OperationID,
		Data:                rawDataPtr,
		Status:              mongoJob.Status,
		SendingAttemptIndex: mongoJob.SendingAttemptIndex,
		SendingError:        mongoJob.SendingError,
		PipelineID:          mongoJob.PipelineID,
		StartAt:             mongoJob.StartAt,
		ProcessedDate:       mongoJob.ProcessedDate,
	}
}

func ModelJobToMongoJob(modelJob Job) entities.Job {
	var bsonDataPtr *bson.Raw
	if modelJob.Data != nil {
		var dataMap map[string]interface{}
		json.Unmarshal([]byte(*modelJob.Data), &dataMap)

		_, data, _ := bson.MarshalValue(dataMap)

		bsonData := bson.Raw(data)
		bsonDataPtr = &bsonData
	}

	return entities.Job{
		ID:                  modelJob.ID,
		OperationID:         modelJob.OperationID,
		Data:                bsonDataPtr,
		Status:              modelJob.Status,
		SendingAttemptIndex: modelJob.SendingAttemptIndex,
		SendingError:        modelJob.SendingError,
		PipelineID:          modelJob.PipelineID,
		StartAt:             modelJob.StartAt,
		ProcessedDate:       modelJob.ProcessedDate,
	}
}
