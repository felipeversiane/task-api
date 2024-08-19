package cache

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     getConnectionString(),
		PoolSize: 100,
	})

	return nil
}

func getConnectionString() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	if password != "" {
		return fmt.Sprintf("%s:%s@%s:%s", password, host, host, port)
	}
	return fmt.Sprintf("%s:%s", host, port)
}
