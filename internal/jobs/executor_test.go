package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExecuteSleep(t *testing.T) {
	job := &Job{
		ID: uuid.New(),
		Type: "sleep",
		Payload: []byte(`{"seconds":1}`),
	}

	start := time.Now()

	result, err := Execute(job)

	duration := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result["message"] != "slept for 1 seconds" {
		t.Fatalf("unexpected result: %v", result)
	}

	if duration < time.Second {
		t.Fatalf("expected job to take at least 1 second, took %v", duration)
	}
}