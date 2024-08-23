package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	domain "github.com/felipeversiane/task-api/internal"
	"github.com/felipeversiane/task-api/internal/rest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TaskRepository struct {
	Database *pgxpool.Pool
	Cache    *redis.Client
}

func NewTaskRepository(database *pgxpool.Pool, cache *redis.Client) TaskRepository {
	return TaskRepository{
		Database: database,
		Cache:    cache,
	}
}

func (r *TaskRepository) Insert(ctx context.Context, task domain.Task) (*TaskResponse, *rest.RestError) {
	nameKey := fmt.Sprintf("task:name:%s", task.Name)

	if v, err := r.Cache.Get(ctx, nameKey).Result(); err == nil && v != "" {
		return nil, rest.NewBadRequestError(fmt.Sprintf("task with name %s already exists", task.Name))
	}

	query := `INSERT INTO tasks (id, name, description, situation, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, name, description, situation, created_at, updated_at`

	var taskResponse TaskResponse
	err := r.Database.QueryRow(ctx, query,
		task.ID, task.Name, task.Description, task.Situation, task.CreatedAt, task.UpdatedAt).
		Scan(&taskResponse.ID, &taskResponse.Name, &taskResponse.Description,
			&taskResponse.Situation, &taskResponse.CreatedAt, &taskResponse.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, rest.NewBadRequestError(fmt.Sprintf("task with name %s already exists", task.Name))
		}
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}

	taskJSON, err := json.Marshal(taskResponse)
	if err != nil {
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}
	_, err = r.Cache.Set(ctx, taskResponse.ID.String(), taskJSON, 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	_, err = r.Cache.Set(ctx, nameKey, taskResponse.ID.String(), 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task name: %v", err))
	}

	return &taskResponse, nil
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, task domain.Task) (*TaskResponse, *rest.RestError) {
	nameKey := fmt.Sprintf("task:name:%s", task.Name)

	existingID, err := r.Cache.Get(ctx, nameKey).Result()
	if err == nil && existingID != "" && existingID != id.String() {
		return nil, rest.NewBadRequestError(fmt.Sprintf("task with name %s already exists", task.Name))
	}

	query := `UPDATE tasks SET name = $1, description = $2, situation = $3, updated_at = $4
	          WHERE id = $5
	          RETURNING id, name, description, situation, created_at, updated_at`

	var taskResponse TaskResponse
	err = r.Database.QueryRow(ctx, query, task.Name, task.Description, task.Situation, task.UpdatedAt, id).
		Scan(&taskResponse.ID, &taskResponse.Name, &taskResponse.Description, &taskResponse.Situation,
			&taskResponse.CreatedAt, &taskResponse.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, rest.NewBadRequestError(fmt.Sprintf("task with name %s already exists", task.Name))
		}
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}

	taskJSON, err := json.Marshal(taskResponse)
	if err != nil {
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}
	_, err = r.Cache.Set(ctx, id.String(), taskJSON, 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	_, err = r.Cache.Set(ctx, nameKey, id.String(), 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task name: %v", err))
	}

	return &taskResponse, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) *rest.RestError {
	taskJSON, err := r.Cache.Get(ctx, id.String()).Result()
	if err != nil {
		return rest.NewNotFoundError(fmt.Sprintf("task with ID %s not found", id))
	}

	var task TaskResponse
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		return rest.NewInternalServerError(fmt.Sprintf("Failed to unmarshal task: %s", err))
	}

	if _, err := r.Cache.Del(ctx, id.String()).Result(); err != nil {
		slog.Error(fmt.Sprintf("Failed to delete task from cache: %v", err))
	}

	nameKey := fmt.Sprintf("task:name:%s", task.Name)
	if _, err := r.Cache.Del(ctx, nameKey).Result(); err != nil {
		slog.Error(fmt.Sprintf("Failed to delete task name from cache: %v", err))
	}

	query := `DELETE FROM tasks WHERE id = $1 RETURNING id`
	var deletedID uuid.UUID
	if err := r.Database.QueryRow(ctx, query, id).Scan(&deletedID); err != nil {
		if err == pgx.ErrNoRows {
			return rest.NewNotFoundError(fmt.Sprintf("task with ID %s not found", id))
		}
		return rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*TaskResponse, *rest.RestError) {
	taskJSON, err := r.Cache.Get(ctx, id.String()).Result()
	if err == nil {
		var task TaskResponse
		if err := json.Unmarshal([]byte(taskJSON), &task); err == nil {
			return &task, nil
		}
	}

	query := `SELECT id, name, description, situation, created_at, updated_at FROM tasks WHERE id = $1`
	row := r.Database.QueryRow(ctx, query, id)
	var task TaskResponse
	if err := row.Scan(&task.ID, &task.Name, &task.Description, &task.Situation, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, rest.NewNotFoundError(fmt.Sprintf("task with ID %s not found", id))
		}
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}
	_, err = r.Cache.Set(ctx, id.String(), string(taskBytes), 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	return &task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]TaskResponse, *rest.RestError) {
	rows, err := r.Database.Query(ctx, `SELECT id, name, description, situation, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}
	defer rows.Close()

	var tasks []TaskResponse
	for rows.Next() {
		var task TaskResponse
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Situation, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, rest.NewInternalServerError(fmt.Sprintf("%s", err))
	}

	return tasks, nil
}
