package task

import (
	"fmt"
	"strings"
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

func (req *TaskRequest) Validate() error {
	var missingFields []string
	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if req.Description == "" {
		missingFields = append(missingFields, "description")
	}
	if req.Situation == "" {
		missingFields = append(missingFields, "situation")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}
	return nil
}

func (req *UpdateTaskRequest) Validate() error {
	var missingFields []string
	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if req.Description == "" {
		missingFields = append(missingFields, "description")
	}
	if req.Situation == "" {
		missingFields = append(missingFields, "situation")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}
	return nil
}
