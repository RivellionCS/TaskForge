package jobs

import (
	"encoding/json"
	"fmt"
	"time"
)

func Execute(job *Job) error {
	switch job.Type {
	case "sleep":
		var payload struct {
			Seconds int `json:"seconds"`
		}

		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return fmt.Errorf("decode sleep payload: %w", err)
		}

		if payload.Seconds < 0 {
			return fmt.Errorf("seconds cannot be negative")
		}

		time.Sleep(time.Duration(payload.Seconds) * time.Second)

		return nil

	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}