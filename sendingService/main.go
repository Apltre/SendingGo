package main

import (
	"encoding/json"
	"log"
	"reflect"
	"sendingService/appconfig"
	"sendingService/controllers"
	"sendingService/models"
	"strconv"
	"time"

	"github.com/streadway/amqp"
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
	_, err := channel.QueueDeclare("sending_queue_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_queue_"+pipelineID+"\"")

	_, err = channel.QueueDeclare("sending_service_"+pipelineID, true, false, false, false, nil)
	failOnError(err, "Can't create queue: \"sending_service_"+pipelineID+"\"")
}

func publishToSendingQueueFromChannel(rabbitChannel *ampq.Channel, dataToSend <-chan []byte, pipelineID string) {
	for {
		dataToPublish := <-dataToSend
		err := rabbitChannel.Publish(
			"",
			"sending_queue_"+pipelineID,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        dataToPublish,
			})

		if err != nil {
			log.Println(err.Error() + " Publish failed for data: " + string(dataToPublish))
		}
	}
}

func handleJobs(channel *ampq.Channel, pipelineID string) {
	sendingQueue, err := channel.Consume(
		"sending_service_"+pipelineID,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Can't connect to sending queue")

	sendingQueuePublishingChannel := make(chan []byte)

	go publishToSendingQueueFromChannel(channel, sendingQueuePublishingChannel, pipelineID)

	for message := range sendingQueue {
		go func(message ampq.Delivery) {
			job := models.Job{}
			err := json.Unmarshal(message.Body, &job)

			if err != nil {
				log.Printf("Unmarshal error job:%s\n", job.ID)
				message.Ack(false)
				return
			}

			controllerData := controllers.ControllersMapping[job.OperationID]
			inputs := []reflect.Value{reflect.ValueOf(job.Data)}
			reflectionResults := reflect.ValueOf(controllerData.ControllerRef).MethodByName(controllerData.MethodName).Call(inputs)
			methodResult := reflectionResults[0].Interface().(*models.SendingError)

			if methodResult == nil {
				job.Status = 2
			} else {
				job.Status = methodResult.ErrorType
				job.SendingError = methodResult.Message
			}

			currentTime := time.Now()
			job.ProcessedDate = &currentTime
			jobJson, _ := json.Marshal(job)
			sendingQueuePublishingChannel <- jobJson

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
