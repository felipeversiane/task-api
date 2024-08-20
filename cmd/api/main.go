package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/felipeversiane/task-api/internal/cache"
	"github.com/felipeversiane/task-api/internal/database"
	"github.com/felipeversiane/task-api/internal/log"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	slog.Info("Configuring logs")
	log.Configure()
	slog.Info("Successfully configured logs")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("Configuring database")
	if err := database.Connect(ctx); err != nil {
		panic(err)
	}
	defer database.Close()
	slog.Info("Successfully configured database")

	slog.Info("Configuring cache")
	if err := cache.Connect(); err != nil {
		panic(err)
	}
	slog.Info("Successfully configured cache")

	slog.Info("Running API")
	go http.ListenAndServe(port, nil)

}
