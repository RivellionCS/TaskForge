package jobs

import (
	"context"
	"testing"

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

func TestRepositoryCreate(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	repository := NewRepository(db)

	jobID, err := repository.Create(
		ctx,
		"sleep",
		map[string]any{
			"seconds": 5,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if jobID == uuid.Nil {
		t.Fatal("expected a job ID")
	}

	job, err := repository.GetByID(ctx, jobID)
	if err != nil {
		t.Fatalf("failed to get created job: %v", err)
	}

	if job.Type != "sleep" {
		t.Fatalf("expected type sleep, got %s", job.Type)
	}

	if job.Status != "pending" {
		t.Fatalf("expected status pending, got %s", job.Status)
	}
}

func TestRepositoryMarkRunning(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	repository := NewRepository(db)

	jobID, err := repository.Create(
		ctx,
		"sleep",
		map[string]any{
			"seconds": 5,
		},
	)

	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	err = repository.MarkRunning(ctx, jobID)
	if err != nil {
		t.Fatalf("failed to mark job as running: %v", err)
	}

	job, err := repository.GetByID(ctx, jobID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if job.Status != "running" {
		t.Fatalf("expected status running, got %s", job.Status)
	}

	if job.StartedAt == nil {
		t.Fatal("expected started_at to be set")
	}
}