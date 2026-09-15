package main

import (
	"context"
	"log"
	"net/http"

	"github.com/RivellionCS/TaskForge/internal/api"
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

	jobRepository := jobs.NewRepository(db)
	jobHandler := api.NewJobHandler(jobRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs", jobHandler.CreateJob)
	mux.HandleFunc("GET /jobs/{id}", jobHandler.GetJob)

	log.Println("TaskForge API listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}