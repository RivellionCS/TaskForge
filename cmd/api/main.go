package main

import (
	"context"
	"log"
	"net/http"

	"github.com/RivellionCS/TaskForge/internal/api"
	"github.com/RivellionCS/TaskForge/internal/database"
	"github.com/RivellionCS/TaskForge/internal/jobs"
)

func main() {
	ctx := context.Background()

	db, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jobRepository := jobs.NewRepository(db)
	jobHandler := api.NewJobHandler(jobRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs", jobHandler.CreateJob)

	log.Println("TaskForge API listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")
}