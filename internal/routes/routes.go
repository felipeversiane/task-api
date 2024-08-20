package routes

import (
	"net/http"

	"github.com/felipeversiane/task-api/internal/task"
)

func SetupRoutes(mux *http.ServeMux) {

	task.TasksRouter(mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
