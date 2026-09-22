package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	jobType string,
	payload map[string]any,
) (uuid.UUID, error) {
	id := uuid.New()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("marshal payload: %w", err)
	}

	_, err = r.db.Exec(
		ctx,
		`
		INSERT INTO jobs (id, type, status, payload)
		VALUES ($1, $2, $3, $4)
		`,
		id,
		jobType,
		"pending",
		payloadJSON,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert job: %w", err)
	}

	return id, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Job, error) {
	var job Job

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, type, status, payload, result, attempts,
		       created_at, started_at, completed_at
		FROM jobs
		WHERE id = $1
		`,
		id,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Status,
		&job.Payload,
		&job.Result,
		&job.Attempts,
		&job.CreatedAt,
		&job.StartedAt,
		&job.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}

	return &job, nil
}

func (r *Repository) MarkRunning(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE jobs
		SET status = $1,
			started_at = NOW()
		WHERE id = $2
		`,
		"running",
		id,
	)

	if err != nil {
		return fmt.Errorf("mark job running: %w", err)
	}

	return nil
}

func (r *Repository) MarkCompleted(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE jobs
		SET status = $1,
			completed_at = NOW()
		WHERE id = $2
		`,
		"completed",
		id,
	)

	if err != nil {
		return fmt.Errorf("mark job completed: %w", err)
	}

	return nil
}

func (r *Repository) IncrementAttempts(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	var attempts int

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE jobs
		SET attempts = attempts + 1
		WHERE id = $1
		RETURNING attempts
		`,
		id,
	).Scan(&attempts)

	if err != nil {
		return 0, fmt.Errorf("increment job attempts: %w", err)
	}

	return attempts, nil
}

func (r *Repository) MarkFailed(
	ctx context.Context,
	id uuid.UUID,
	errMessage string,
) error {
	resultJSON, err := json.Marshal(map[string]any{
		"error": errMessage,
	})

	if err != nil {
		return fmt.Errorf("marshal failure result: %w", err)
	}

	_, err = r.db.Exec(
		ctx,
		`
		UPDATE jobs
		SET status = $1,
			result = $2,
			completed_at = NOW()
		WHERE id = $3
		`,
		"failed",
		resultJSON,
		id,
	)

	if err != nil {
		return fmt.Errorf("mark job failed: %w", err)
	}

	return nil
}

func (r *Repository) SetResult(
	ctx context.Context,
	id uuid.UUID,
	result map[string]any,
) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal job result: %w", err)
	}

	_, err = r.db.Exec(
		ctx,
		`
		UPDATE jobs
		SET result = $1
		WHERE id = $2
		`,
		resultJSON,
		id,
	)

	if err != nil {
		return fmt.Errorf("set job result: %w", err)
	}

	return nil
}