package repositories

import (
	"jobcreator/data/dbcontexts"
	"jobcreator/data/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	_, err := jobsRepository.mongoDbContext.JobsCollection.InsertOne(ctx, job)
	return err
}

func (jobsRepository *JobsRepository) CreateJobsCollectionStatusIndex() error {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"status": 1,
		},
		Options: nil,
	}
	ctx, _ := jobsRepository.mongoDbContext.CreateCommandContext()

	_, err := jobsRepository.mongoDbContext.JobsCollection.Indexes().CreateOne(ctx, mod)
	return err
}

func (jobsRepository *JobsRepository) CreateJobsCollectionStartAtIndex() error {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"startAt": -1,
		},
		Options: nil,
	}
	ctx, _ := jobsRepository.mongoDbContext.CreateCommandContext()

	_, err := jobsRepository.mongoDbContext.JobsCollection.Indexes().CreateOne(ctx, mod)
	return err
}
