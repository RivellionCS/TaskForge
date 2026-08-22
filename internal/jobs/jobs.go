package jobs

import "github.com/google/uuid"

type Job struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Payload     []byte    `json:"payload"`
	Result      []byte    `json:"result,omitempty"`
	Attempts    int       `json:"attempts"`
	CreatedAt   string    `json:"created_at"`
	StartedAt   *string   `json:"started_at,omitempty"`
	CompletedAt *string   `json:"completed_at,omitempty"`
}