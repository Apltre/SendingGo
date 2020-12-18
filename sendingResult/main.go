package main

import (
	"encoding/json"
	"log"
	"reflect"
	"sendingService/appconfig"
	"sendingService/controllers"
	"sendingService/models"
	"strconv"

	ampq "github.com/streadway/amqp"
)

var (
	config appconfig.AppConfig
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s\n", msg, err)
	}
}

func initializeQueues(channel *ampq.Channel, pipelineID string) {
	_, err := channel.QueueDeclare("sending_results_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_results_"+pipelineID+"\"")
}

func handleJobs(channel *ampq.Channel, pipelineID string) {
	resultsQueue, err := channel.Consume(
		"sending_results_"+pipelineID,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Can't connect to sending results queue")

	for message := range resultsQueue {
		go func(message ampq.Delivery) {
			job := models.Job{}
			err := json.Unmarshal(message.Body, &job)

			if err != nil {
				log.Printf("Unmarshal error job: %s\n", job.ID)
				message.Ack(false)
				return
			}

			var methodNameEnding string
			switch job.Status {
			case 2:
				methodNameEnding = "Success"
			case -1:
				methodNameEnding = "Failure"
			default:
				log.Printf("Wrong job status/ Id = %s\n", job.ID.String())
				message.Ack(false)
				return
			}

			controllerData := controllers.ControllersMapping[job.OperationID]
			inputs := []reflect.Value{reflect.ValueOf(job.Data)}
			reflectionResults := reflect.ValueOf(controllerData.ControllerRef).MethodByName(controllerData.MethodName + methodNameEnding).Call(inputs)
			errorResult := reflectionResults[0].Interface()

			if errorResult != nil {
				log.Println("Job " + job.ID.String() + " result handling fail. Error: " + errorResult.(error).Error())
			}

			message.Ack(false)
		}(message)
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

	go handleJobs(ampqChannel, pipelineID)

	wait := make(chan bool)
	<-wait
}
