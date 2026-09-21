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

		if err := repository.MarkRunning(ctx, job.ID); err != nil {
			log.Printf("Failed to update job %s status: %v", job.ID, err)
			message.Nack(false, true)
			continue
		}

		log.Printf(
			"Job started: id=%s type=%s status=running",
			job.ID,
			job.Type,

		)

		result, err := jobs.Execute(job)
		if err != nil {
			attempts, attemptsErr := repository.IncrementAttempts(ctx, job.ID)
			if attemptsErr != nil {
				log.Printf(
					"Failed to increment attempts for job %s: %v",
					job.ID,
					attemptsErr,
				)
				message.Nack(false, true)
				continue
			}

			log.Printf(
				"Job %s failed: %v (attempt %d)",
				job.ID,
				err,
				attempts,
			)

			if attempts >= 3 {
				if err := repository.MarkFailed(ctx, job.ID); err != nil {
					log.Printf(
						"Failed to mark job %s as failed: %v",
						job.ID,
						err,
					)
					message.Nack(false, true)
					continue
				}

				log.Printf("Job permanently failed: id=%s", job.ID)
				message.Ack(false)
				continue
			}

			message.Nack(false, true)
			continue
		}

		if err := repository.SetResult(ctx, job.ID, result); err != nil {
			log.Printf("Failed to save result for job %s: %v", job.ID, err)
			message.Nack(false, true)
			continue
		}

		if err := repository.MarkCompleted(ctx, job.ID); err != nil {
			log.Printf("Failed to mark job %s as completed: %v", job.ID, err)
			message.Nack(false, true)
			continue
		}

		log.Printf("Job completed: id=%s", job.ID)

		message.Ack(false)
	}
}