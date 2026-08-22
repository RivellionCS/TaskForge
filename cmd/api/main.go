package main

import (
	"context"
	"log"

	"github.com/RivellionCS/TaskForge/internal/database"
)

func main() {
	ctx := context.Background()

	db, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")
}