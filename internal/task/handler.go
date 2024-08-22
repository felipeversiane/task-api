package task

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/felipeversiane/task-api/internal/rest"
	"github.com/google/uuid"
)

type TaskHandler struct {
	Service TaskService
}

func NewTaskHandler(service TaskService) TaskHandler {
	return TaskHandler{
		Service: service,
	}
}

func (h *TaskHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr := rest.NewBadRequestError("invalid request payload")
		respondWithJSON(w, httpErr.Code, httpErr)
		return
	}

	resp, err := h.Service.CreateTask(ctx, req)
	if err != nil {
		respondWithJSON(w, err.Code, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, parseErr := extractIDFromPath(r.URL.Path)
	if parseErr != nil {
		httpErr := rest.NewBadRequestError("invalid task ID")
		respondWithJSON(w, httpErr.Code, httpErr)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr := rest.NewBadRequestError("invalid request payload")
		respondWithJSON(w, httpErr.Code, httpErr)
		return
	}

	resp, err := h.Service.UpdateTask(ctx, id, req)
	if err != nil {
		respondWithJSON(w, err.Code, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, parseErr := extractIDFromPath(r.URL.Path)
	if parseErr != nil {
		httpErr := rest.NewBadRequestError("invalid task ID")
		respondWithJSON(w, httpErr.Code, httpErr)
		return
	}

	if err := h.Service.DeleteTask(ctx, id); err != nil {
		respondWithJSON(w, err.Code, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, parseErr := extractIDFromPath(r.URL.Path)
	if parseErr != nil {
		httpErr := rest.NewBadRequestError("invalid task ID")
		respondWithJSON(w, httpErr.Code, httpErr)
		return
	}

	resp, err := h.Service.GetTaskByID(ctx, id)
	if err != nil {
		respondWithJSON(w, err.Code, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.Service.GetAllTasks(ctx)
	if err != nil {
		respondWithJSON(w, err.Code, err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func extractIDFromPath(path string) (uuid.UUID, error) {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	return uuid.Parse(idStr)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
