package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RivellionCS/TaskForge/internal/jobs"
	"github.com/RivellionCS/TaskForge/internal/queue"
	"github.com/google/uuid"
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

	var response struct {
		ID string `json:"id"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == "" {
		t.Fatal("expected job ID in response")
	}

	if response.Status != "pending" {
		t.Fatalf("expected status pending, got %s", response.Status)
	}

	jobID, err := uuid.Parse(response.ID)
	if err != nil {
		t.Fatalf("invalid job ID returned by API: %v", err)
	}

	job, err := repository.GetByID(request.Context(), jobID)
	if err != nil {
		t.Fatalf("failed toget created job: %v", err)
	}

	if job.Type != "sleep" {
		t.Fatalf("expected job type sleep, got %s", job.Type)
	}

	if job.Status != "pending" {
		t.Fatalf("expected job status pending, got %s", job.Status)
	}

	if string(job.Payload) != `{"seconds": 5}` {
		t.Fatalf("unexpected job payload: %s", job.Payload)
	}
}

func TestGetJob(t *testing.T) {
	db := newTestDB(t)

	repository := jobs.NewRepository(db)

	jobID, err := repository.Create(
		context.Background(),
		"sleep",
		map[string]any{
			"seconds": 5,
		},
	)

	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	handler := NewJobHandler(repository, nil)

	request := httptest.NewRequest(
		http.MethodGet,
		"/jobs/"+jobID.String(),
		nil,
	)

	request.SetPathValue("id", jobID.String())

	recorder := httptest.NewRecorder()

	handler.GetJob(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response jobs.Job
	
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != jobID {
		t.Fatalf("expected job ID %s, got %s", jobID, response.ID)
	}

	if response.Type != "sleep" {
		t.Fatalf("expected job type sleep, got %s", response.Type)
	}

	if response.Status != "pending" {
		t.Fatalf("expected job status pending, got %s", response.Status)
	}

	if string(response.Payload) != `{"seconds": 5}` {
		t.Fatalf("unexpected payload: %s", response.Payload)
	}
}

func TestCreateJobMissingType(t *testing.T) {
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

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}