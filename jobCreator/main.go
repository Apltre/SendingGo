package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"jobcreator/data/dbcontexts"
	"jobcreator/data/entities"
	"jobcreator/data/repositories"

	"jobcreator/appconfig"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	config       appconfig.AppConfig
	mongoContext *dbcontexts.MongoDbContext
)

func returnInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	var requestObject JobDto
	err := json.NewDecoder(r.Body).Decode(&requestObject)
	if err != nil {
		log.Println(err)
		returnInternalServerError(w, err)
	}

	startAt := time.Now()

	if requestObject.StartAt != nil {
		startAt = *requestObject.StartAt
	}

	job := entities.Job{
		ID:          primitive.NewObjectID(),
		Status:      requestObject.Status,
		OperationID: requestObject.OperationID,
		Data:        requestObject.Data,
		PipelineID:  requestObject.PipelineID,
		StartAt:     startAt,
	}
	jobsRepository := repositories.NewJobsRepository(mongoContext)
	err = jobsRepository.InsertJob(job)
	if err != nil {
		log.Println(err)
		returnInternalServerError(w, err)
	}
}

func main() {
	config = appconfig.NewConfig()
	mongoContext = dbcontexts.NewMongoDbContext(config)
	jobsRepository := repositories.NewJobsRepository(mongoContext)
	err := jobsRepository.CreateJobsCollectionStatusIndex()

	if err != nil {
		log.Fatalln(err)
	}

	err = jobsRepository.CreateJobsCollectionStartAtIndex()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Initializing...")
	http.HandleFunc("/jobs", jobHandler)

	err = http.ListenAndServe(":"+strconv.Itoa(config.APIConfig.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
