package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sendingQueue/appconfig"
	"sendingQueue/data/dbcontexts"
	"sendingQueue/data/repositories"
	"sendingQueue/models"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	ampq "github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

var (
	config       appconfig.AppConfig
	mongoContext *dbcontexts.MongoDbContext
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s\n", msg, err)
	}
}

func initializeQueues(channel *ampq.Channel, pipelineID string) {

	_, err := channel.QueueDeclare("sending_queue_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_queue_"+pipelineID+"\"")

	_, err = channel.QueueDeclare("sending_results_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_results_"+pipelineID+"\"")

	_, err = channel.QueueDeclare("sending_service_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_service_"+pipelineID+"\"")
}

func handleNewJobs(channel *ampq.Channel, pipelineID int) {
	mongoContext = dbcontexts.NewMongoDbContext(config)
	rep := repositories.NewJobsRepository(mongoContext)

	for {
		items, err := rep.GetNewJobsPack(10, pipelineID)
		if err != nil {
			log.Println(err)
		}

		if len(items) == 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, item := range items {
			item.Status = 1

			jobModel := models.MongoJobToModel(item)
			json, err := json.Marshal(jobModel)

			if err != nil {
				fmt.Println(err)
				continue
			}

			err = channel.Publish(
				"",
				"sending_service_"+strconv.Itoa(pipelineID),
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        json,
				})

			if err != nil {
				fmt.Println(err)
				continue
			}

			err = rep.UpdateJob(item)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func handleJobs(channel *ampq.Channel, pipelineID string) {
	sendingQueue, err := channel.Consume(
		"sending_queue_"+pipelineID,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Can't connect to sending queue.")

	for message := range sendingQueue {
		job := models.Job{}

		err := json.Unmarshal(message.Body, &job)

		if err != nil {
			log.Printf("Error parsing rabbitMq message to json: %s\n", err.Error())
			continue
		}

		mongoContext = dbcontexts.NewMongoDbContext(config)
		rep := repositories.NewJobsRepository(mongoContext)

		currentTime := time.Now()
		job.ProcessedDate = &currentTime

		mongojob := models.ModelJobToMongoJob(job)
		err = rep.UpdateJob(mongojob)
		if err != nil {
			log.Printf("Job update failed ID = %s %s\n", job.ID.Hex(), err.Error())
			continue
		}

		switch job.Status {
		case -1, 2:
			err = channel.Publish(
				"",
				"sending_results_"+pipelineID,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        message.Body,
				})

			if err != nil {
				log.Println(err)
				continue
			}

		case -2:
			mongojob.StartAt = currentTime.Add(5 * time.Minute)
			mongojob.SendingAttemptIndex++
			mongojob.Status = 0
			mongojob.SendingError = nil
			mongojob.ProcessedDate = nil
			mongojob.ID = primitive.NewObjectID()
			err = rep.InsertJob(mongojob)

			if err != nil {
				log.Printf("New job insertion failed. Failed job ID = %s, Error: %s\n", job.ID.String(), err.Error())
				continue
			}
		}

		message.Ack(true)
	}
}

func main() {
	config = appconfig.NewConfig()
	pipelineID := strconv.Itoa(config.PipelineID)
	conn, err := ampq.Dial(config.RabbitMq.ConnectionString)
	defer conn.Close()

	failOnError(err, "Can't connect to RabbitMQ")

	var ampqChannel *ampq.Channel
	ampqChannel, err = conn.Channel()
	defer ampqChannel.Close()

	failOnError(err, "Can't create ampq channel")

	initializeQueues(ampqChannel, pipelineID)

	log.Println("Rabbit queues initialized")

	go handleNewJobs(ampqChannel, config.PipelineID)
	go handleJobs(ampqChannel, pipelineID)

	wait := make(chan bool)
	<-wait
}
