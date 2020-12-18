package repositories

import (
	"context"
	"errors"
	"sendingQueue/data/dbcontexts"
	"sendingQueue/data/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobsRepository struct {
	mongoDbContext *dbcontexts.MongoDbContext
}

func (jobsRepository *JobsRepository) configJobsRepository(mongoDbContext *dbcontexts.MongoDbContext) {
	jobsRepository.mongoDbContext = mongoDbContext
}

func NewJobsRepository(mongoContext *dbcontexts.MongoDbContext) *JobsRepository {
	var rep = &JobsRepository{}
	rep.configJobsRepository(mongoContext)
	return rep
}

func (jobsRepository *JobsRepository) InsertJob(job entities.Job) error {
	ctx, _ := jobsRepository.mongoDbContext.CreateCommandContext()
	data, err := bson.Marshal(job)
	_, err = jobsRepository.mongoDbContext.JobsCollection.InsertOne(ctx, data)
	return err
}

func (jobsRepository *JobsRepository) UpdateJob(job entities.Job) error {
	ctx, _ := jobsRepository.mongoDbContext.CreateCommandContext()

	filter := bson.D{{"_id", job.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"status", job.Status},
			{"processedDate", job.ProcessedDate},
			{"sendingAttemptIndex", job.SendingAttemptIndex},
			{"sendingError", job.SendingError},
			{"startAt", job.StartAt},
		}},
	}
	_, err := jobsRepository.mongoDbContext.JobsCollection.UpdateOne(ctx, filter, update)
	return err
}

func (jobsRepository *JobsRepository) GetNewJobsPack(size int, pipelineID int) ([]entities.Job, error) {
	ctx, _ := jobsRepository.mongoDbContext.CreateCommandContext()
	if size < 0 {
		return nil, errors.New("jobs batch size can't be less than zero")
	}
	options := options.Find()
	options.SetLimit(int64(size))
	options.SetSort(bson.M{"_id": 1})
	filter := bson.D{{"status", 0}, {"pipelineId", pipelineID}, {"startAt", bson.M{"$lte": time.Now()}}}
	var result []entities.Job
	cursor, err := jobsRepository.mongoDbContext.JobsCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {

		var job entities.Job
		err := cursor.Decode(&job)
		if err != nil {
			return nil, err
		}

		result = append(result, job)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
