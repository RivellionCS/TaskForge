package main

import (
	"context"
	"log"

	"github.com/RivellionCS/TaskForge/internal/database"
	"github.com/RivellionCS/TaskForge/internal/jobs"
	"github.com/RivellionCS/TaskForge/internal/queue"
	"github.com/google/uuid"
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
		jobID, err := uuid.Parse(string(message.Body))
		if err != nil {
			log.Printf("Invalid job ID: %s", message.Body)
			message.Nack(false, false)
			continue
		}

		job, err := repository.GetByID(ctx, jobID)
		if err != nil {
			log.Printf("Failed to get job %s: %v", jobID, err)
			message.Nack(false, true)
			continue
		}

		log.Printf(
			"Recieved job: id=%s type=%s status=%s",
			job.ID,
			job.Type,
			job.Status,
		)
	}
}