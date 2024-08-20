package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/felipeversiane/task-api/internal/cache"
	"github.com/felipeversiane/task-api/internal/database"
	"github.com/felipeversiane/task-api/internal/log"
	"github.com/felipeversiane/task-api/internal/routes"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	log.Configure()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := database.Connect(ctx); err != nil {
		panic(err)
	}
	defer database.Close()

	if err := cache.Connect(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	slog.Info(fmt.Sprintf("Server running on port : %s", port))
	http.ListenAndServe(":"+port, mux)
}
