package main

import (
	"context"
	"log"

	"github.com/RivellionCS/TaskForge/internal/database"
	"github.com/RivellionCS/TaskForge/internal/jobs"
	"github.com/RivellionCS/TaskForge/internal/queue"
)

func main() {
	ctx := context.Background()

	db, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rabbitmq, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	if err := rabbitmq.DeclareQueue(); err != nil {
		log.Fatal(err)
	}

	repository := jobs.NewRepository(db)

	messages, err := rabbitmq.ConsumeJobs()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TaskForge worker waiting for jobs...")
	
	for message := range messages {
		log.Printf("Received job: %s", message.Body)
	}
}