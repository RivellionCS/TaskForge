package main

import (
	"log"

	"github.com/RivellionCS/TaskForge/internal/queue"
)

func main() {
	rabbitmq, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	if err := rabbitmq.DeclareQueue(); err != nil {
		log.Fatal(err)
	}

	messages, err := rabbitmq.ConsumeJobs()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TaskForge worker waiting for jobs...")
	
	for message := range messages {
		log.Printf("Received job: %s", message.Body)
	}
}