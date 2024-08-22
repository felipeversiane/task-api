package task

import (
	"context"

	domain "github.com/felipeversiane/task-api/internal"
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

func (s *TaskService) CreateTask(ctx context.Context, domain domain.Task) (*TaskResponse, error) {
	task, err := s.Repository.Insert(ctx, domain)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, domain domain.Task) (*TaskResponse, error) {
	task, err := s.Repository.Update(ctx, id, domain)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	_, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.Repository.Delete(ctx, id)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error) {
	task, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]TaskResponse, error) {
	tasks, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
