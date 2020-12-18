package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID
	OperationID         int
	Data                *json.RawMessage
	Status              int
	SendingAttemptIndex int
	SendingError        *string
	PipelineID          int
	StartAt             time.Time
	ProcessedDate       *time.Time
}
