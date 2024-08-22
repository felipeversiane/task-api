package task

import (
	"context"

	"github.com/felipeversiane/task-api/internal/rest"
	"github.com/google/uuid"
)

type TaskService struct {
	Repository TaskRepository
}

func NewTaskService(repository TaskRepository) TaskService {
	return TaskService{
		Repository: repository,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req TaskRequest) (*TaskResponse, *rest.RestError) {

	if err := req.Validate(); err != nil {
		return nil, rest.NewBadRequestError(err.Error())
	}

	domain := RequestToDomainTask(req)
	if err := domain.ValidateFields(); err != nil {
		return nil, rest.NewBadRequestError(err.Error())
	}

	task, err := s.Repository.Insert(ctx, domain)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, *rest.RestError) {

	if err := req.Validate(); err != nil {
		rest.NewBadRequestError(err.Error())
	}

	domain := RequestToUpdateDomainTask(req)
	if err := domain.ValidateFields(); err != nil {
		rest.NewBadRequestError(err.Error())
	}

	task, err := s.Repository.Update(ctx, id, domain)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) *rest.RestError {
	_, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.Repository.Delete(ctx, id)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*TaskResponse, *rest.RestError) {
	task, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]TaskResponse, *rest.RestError) {
	tasks, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
