package task

import (
	"time"

	domain "github.com/felipeversiane/task-api/internal"
	"github.com/google/uuid"
)

type TaskRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Situation   domain.Situation `json:"situation"`
}

type UpdateTaskRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Situation   domain.Situation `json:"situation"`
}

type TaskResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Situation   domain.Situation `json:"situation"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func RequestToDomainTask(req TaskRequest) domain.Task {
	return domain.NewTask(
		req.Name,
		req.Description,
		req.Situation,
	)
}

func RequestToUpdateDomainTask(req UpdateTaskRequest) domain.Task {
	return domain.NewUpdateTask(
		req.Name,
		req.Description,
		req.Situation,
	)
}

func DomainToResponseTask(domain domain.Task) TaskResponse {
	return TaskResponse{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		Situation:   domain.Situation,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}
