package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RivellionCS/TaskForge/internal/jobs"
	"github.com/RivellionCS/TaskForge/internal/queue"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://taskforge:taskforge@localhost:5432/taskforge_test",
	)
	if err != nil {
		t.Fatalf("failed to create test database pool: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		t.Fatalf("failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestCreateJob(t *testing.T) {
	db := newTestDB(t)

	repository := jobs.NewRepository(db)

	rabbitmq, err := queue.NewRabbitMQ()
	if err != nil {
		t.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	if err := rabbitmq.DeclareQueue(); err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}

	handler := NewJobHandler(repository, rabbitmq)

	body := `{
		 "type": "sleep",
		 "payload": {
		 	 "seconds": 5
		 }
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/jobs",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.CreateJob(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}
}