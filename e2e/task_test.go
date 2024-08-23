package e2e

import (
	"net/http"
	"testing"

	"github.com/felipeversiane/task-api/internal/task"
	"github.com/google/uuid"
)

func happyData() task.TaskRequest {
	return task.TaskRequest{
		Name:        "Beautiful task.",
		Description: "A beautiful task to do.",
		Situation:   "in progress",
	}
}

func TestInsertTask_ShouldReturnStatusBadRequest_WhenItHasInvalidData(t *testing.T) {
	t.Log("*** Test Insert Task with Invalid Data")

	api := NewApiClient()
	params := []map[string]interface{}{
		nil,
		{},
		{"other": "value"},
		{"name": "Ta", "description": "A beautiful task to do.", "situation": "in progress"},
		{"name": "", "description": "A beautiful task to do.", "situation": "in progress"},
		{"name": "Task 1", "description": "A beautiful task to do.", "situation": "blocked"},
	}
	for _, p := range params {
		resp, err := api.Post("/tasks", p)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()
		assertStatusCode(t, resp, http.StatusBadRequest)
	}

}

func TestInsertDuplicateTask_ShouldReturnStatusBadRequest_DuplicatedNameInDatabase(t *testing.T) {
	t.Log("*** Test Insert Duplicate Task")

	api := NewApiClient()
	task := happyData()

	id := insertTaskSuccessfully(task, t)

	payload := map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
		"situation":   task.Situation,
	}

	resp, err := api.Post("/tasks", payload)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusBadRequest)

	deleteTaskSuccessfully(id, t)
}

func TestDeleteTask_ShouldReturnStatusNotFound_WhenTaskIsNotOnDatabase(t *testing.T) {
	t.Log("*** Test Delete Task when Task is not on Database")

	api := NewApiClient()
	id := uuid.NewString()

	resp, err := api.Delete("/tasks/" + id)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusNotFound)

}

func TestGetTaskByID_ShouldReturnStatusNotFound_WhenTaskIsNotOnDatabase(t *testing.T) {
	t.Log("*** Test Get Task by ID when Task is not on Database")

	api := NewApiClient()
	id := uuid.NewString()

	resp, err := api.Get("/tasks/" + id)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	assertStatusCode(t, resp, http.StatusNotFound)

}

func insertTaskSuccessfully(task task.TaskRequest, t *testing.T) string {
	t.Log("***Test Insert Task Successfully")
	api := NewApiClient()

	payload := map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
		"situation":   task.Situation,
	}

	resp, err := api.Post("/tasks", payload)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusCreated)

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err)
	}

	id := res["id"].(string)
	if id == "" {
		t.Fatal("Invalid ID")
	}
	if res["name"].(string) != task.Name {
		t.Fatal("Invalid Name")
	}
	if res["created_at"].(string) == "0001-01-01T00:00:00Z" {
		t.Fatal("Invalid CreatedAt")
	}

	return id
}

func getTaskByIDSuccessfully(id string, t *testing.T) {
	t.Log("***Test Get Task By ID Successfully")

	task := happyData()
	api := NewApiClient()

	resp, err := api.Get("/tasks/" + id)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusOK)

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err)
	}

	if res["id"].(string) != id {
		t.Fatal("Invalid ID")
	}
	if res["name"].(string) != task.Name {
		t.Fatal("Invalid name")
	}

}

func updateTaskSuccessfully(id string, t *testing.T) {
	t.Log("***Test Update Task Successfully")

	api := NewApiClient()
	task := happyData()

	payload := map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
		"situation":   task.Situation,
	}

	resp, err := api.Put("/tasks/"+id, payload)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusOK)

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err)
	}

	if res["id"].(string) != id {
		t.Fatal("Invalid ID")
	}
	if res["name"].(string) != task.Name {
		t.Fatal("Invalid Name")
	}
	if res["created_at"].(string) == "0001-01-01T00:00:00Z" {
		t.Fatal("Invalid CreatedAt")
	}

}

func deleteTaskSuccessfully(id string, t *testing.T) {
	t.Log("***Test Delete Task Successfully")

	api := NewApiClient()

	resp, err := api.Delete("/tasks/" + id)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertStatusCode(t, resp, http.StatusNoContent)

}

func TestUserFlow(t *testing.T) {
	t.Log("*** Start Task Flow ")

	task := happyData()
	id := insertTaskSuccessfully(task, t)
	getTaskByIDSuccessfully(id, t)
	updateTaskSuccessfully(id, t)
	deleteTaskSuccessfully(id, t)

	t.Log("*** End Task Flow Successfull")
}
