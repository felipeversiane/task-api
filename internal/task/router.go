package task

import (
	"net/http"

	"github.com/felipeversiane/task-api/internal/cache"
	"github.com/felipeversiane/task-api/internal/database"
)

var Handler TaskHandler

func TasksRouter(mux *http.ServeMux) {
	Handler = NewTaskHandler(NewTaskService(NewTaskRepository(database.Connection, cache.Client)))

	mux.HandleFunc("POST /api/v1/tasks", Handler.PostTask)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", Handler.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", Handler.DeleteTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", Handler.GetTaskByID)
	mux.HandleFunc("GET /api/v1/tasks", Handler.GetAllTasks)
}
