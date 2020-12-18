package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Job is a main structure containing all information needed for data sending to outer service persisted in mongo db
type Job struct {
	ID                  primitive.ObjectID `bson:"_id"`
	OperationID         int                `bson:"operationId"`
	Data                *bson.Raw          `bson:"data"`
	Status              int                `bson:"status"`
	SendingAttemptIndex int                `bson:"sendingAttemptIndex"`
	SendingError        *string            `bson:"sendingError"`
	PipelineID          int                `bson:"pipelineId"`
	StartAt             time.Time          `bson:"startAt"`
	ProcessedDate       *time.Time         `bson:"processedDate"`
}
