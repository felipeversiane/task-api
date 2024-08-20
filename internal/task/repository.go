package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	domain "github.com/felipeversiane/task-api/internal"
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

func (r *TaskRepository) Insert(ctx context.Context, task domain.Task) (*TaskResponse, error) {
	nameKey := fmt.Sprintf("task:name:%s", task.Name)
	if v, err := r.Cache.Get(ctx, nameKey).Result(); err == nil && v != "" {
		return nil, fmt.Errorf("task with name %s already exists", task.Name)
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
		return nil, err
	}

	taskJSON, err := json.Marshal(taskResponse)
	if err != nil {
		return nil, err
	}

	_, err = r.Cache.Set(ctx, taskResponse.ID.String(), taskJSON, 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	_, err = r.Cache.Set(ctx, nameKey, taskResponse.ID.String(), 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	return &taskResponse, nil
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, domain domain.Task) (*TaskResponse, error) {
	nameKey := fmt.Sprintf("task:name:%s", domain.Name)
	existingID, err := r.Cache.Get(ctx, nameKey).Result()
	if err == nil && existingID != "" && existingID != id.String() {
		return nil, fmt.Errorf("task with name %s already exists", domain.Name)
	}

	query := `UPDATE tasks SET name = $1, description = $2, situation = $3, updated_at = $4
	          WHERE id = $5
	          RETURNING id, name, description, situation, created_at, updated_at`

	var taskResponse TaskResponse
	err = r.Database.QueryRow(ctx, query, domain.Name, domain.Description, domain.Situation, domain.UpdatedAt, id).
		Scan(&taskResponse.ID, &taskResponse.Name, &taskResponse.Description, &taskResponse.Situation,
			&taskResponse.CreatedAt, &taskResponse.UpdatedAt)
	if err != nil {
		return nil, err
	}

	taskJSON, err := json.Marshal(taskResponse)
	if err != nil {
		return nil, err
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

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.Cache.Del(ctx, id.String()).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete task from cache: %v", err))
	}

	query := `DELETE FROM tasks WHERE id = $1 RETURNING id`
	var deletedID uuid.UUID
	err = r.Database.QueryRow(ctx, query, id).Scan(&deletedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("task with ID %s not found", id)
		}
		return err
	}

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error) {
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
			return nil, fmt.Errorf("task with ID %s not found", id)
		}
		return nil, err
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	_, err = r.Cache.Set(ctx, id.String(), string(taskBytes), 24*time.Hour).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to cache task: %v", err))
	}

	return &task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]TaskResponse, error) {
	rows, err := r.Database.Query(ctx, `SELECT id, name, description, situation, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskResponse
	for rows.Next() {
		var task TaskResponse
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Situation, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
