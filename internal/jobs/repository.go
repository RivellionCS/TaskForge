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