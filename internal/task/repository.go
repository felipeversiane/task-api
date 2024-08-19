package task

import (
	domain "github.com/felipeversiane/task-api/internal"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TaskRepository struct {
	Database *pgxpool.Pool
	Cache    *redis.Client
	ChTasks  chan domain.Task
}

func NewTaskRepository(database *pgxpool.Pool, cache *redis.Client) TaskRepository {
	return TaskRepository{
		Database: database,
		Cache:    cache,
		ChTasks:  make(chan domain.Task),
	}
}
