package jobs

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Payload     []byte    `json:"payload"`
	Result      []byte    `json:"result,omitempty"`
	Attempts    int       `json:"attempts"`
	CreatedAt   time.Time    `json:"created_at"`
	StartedAt   *time.Time   `json:"started_at,omitempty"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}