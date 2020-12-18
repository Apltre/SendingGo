package dbcontexts

import (
	"context"
	"sendingQueue/appconfig"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbContext struct {
	client         *mongo.Client
	commandTimeout int
	JobsCollection *mongo.Collection
	Id             string
}

func (dbContext *MongoDbContext) CreateCommandContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbContext.commandTimeout)*time.Second)
}

func (dbContext *MongoDbContext) configureContext(config appconfig.MongoConfig) {
	var err error
	var client *mongo.Client
	dbContext.commandTimeout = config.CommandTimeout

	ctx, _ := dbContext.CreateCommandContext()
	clientOptions := options.Client().ApplyURI("mongodb://" + config.ConnectionString)
	client, err = mongo.Connect(ctx, clientOptions)
	dbContext.client = client

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	dbContext.Id = uuid.New().String()
	dbContext.JobsCollection = client.Database("Sending").Collection("Jobs")
}

func (dbContext *MongoDbContext) Close() {
	ctx, _ := dbContext.CreateCommandContext()
	dbContext.client.Disconnect(ctx)
}

func NewMongoDbContext(config appconfig.AppConfig) *MongoDbContext {
	var dbContext = &MongoDbContext{}
	dbContext.configureContext(config.MongoConfig)
	return dbContext
}
