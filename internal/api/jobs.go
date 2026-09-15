package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RivellionCS/TaskForge/internal/jobs"
	"github.com/RivellionCS/TaskForge/internal/queue"
	"github.com/google/uuid"
)

type JobHandler struct {
	repository *jobs.Repository
	rabbitmq *queue.RabbitMQ
}

func NewJobHandler(
	repository *jobs.Repository,
	rabbitmq *queue.RabbitMQ,
	) *JobHandler {
	return &JobHandler{
		repository: repository,
		rabbitmq: rabbitmq,
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

	if err := h.rabbitmq.PublishJob(id.String()); err != nil {
		http.Error(w, "failed to publish job", http.StatusInternalServerError)
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

func (h *JobHandler) GetJob( w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := h.repository.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("GetByID error: %v", err)
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}