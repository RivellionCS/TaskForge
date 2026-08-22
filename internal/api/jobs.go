package api

import (
	"encoding/json"
	"net/http"

	"github.com/RivellionCS/TaskForge/internal/jobs"
)

type JobHandler struct {
	repository *jobs.Repository
}

func NewJobHandler(repository *jobs.Repository) *JobHandler {
	return &JobHandler{
		repository: repository,
	}
}

type createJobRequest struct {
	Type 	string		   `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var request createJobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	id, err := h.repository.Create(
		r.Context(),
		request.Type,
		request.Payload,
	)

	if err != nil {
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id": id,
		"status": "pending",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}