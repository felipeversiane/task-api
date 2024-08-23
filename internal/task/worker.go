package task

import (
	"context"
	"errors"
	"time"

	domain "github.com/felipeversiane/task-api/internal"
)

type Method string

const (
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	GET    Method = "GET"
)

var validMethods = map[Method]bool{
	POST:   true,
	PUT:    true,
	DELETE: true,
	GET:    true,
}

func IsValidMethod(m Method) bool {
	return validMethods[m]
}

type Worker struct {
	Service TaskService
	Batch   int
	Timeout time.Duration
	Method  Method
}

func NewWorker(service TaskService, batch int, timeout time.Duration, method Method) (*Worker, error) {
	if !IsValidMethod(method) {
		return nil, errors.New("invalid method value")
	}
	return &Worker{
		Service: service,
		Batch:   batch,
		Timeout: timeout,
		Method:  method,
	}, nil
}

func (w *Worker) RunWorker(ctx context.Context, chTasks chan domain.Task, chExit chan struct{}) {

}
