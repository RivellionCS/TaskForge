package jobs

import (
	"encoding/json"
	"fmt"
	"time"
)

func Execute(job *Job) (map[string]any, error) {
	switch job.Type {
	case "sleep":
		var payload struct {
			Seconds int `json:"seconds"`
		}

		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return nil ,fmt.Errorf("decode sleep payload: %w", err)
		}

		if payload.Seconds < 0 {
			return nil, fmt.Errorf("seconds cannot be negative")
		}

		time.Sleep(time.Duration(payload.Seconds) * time.Second)

		return map[string]any{
			"message": fmt.Sprintf("slept for %d seconds", payload.Seconds),
		}, nil

	default:
		return nil, fmt.Errorf("unknown job type: %s", job.Type)
	}
}