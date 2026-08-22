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