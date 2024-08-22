package rest

import (
	"net/http"
)

type RestError struct {
	Message string `json:"message"`
	Err     string `json:"error"`
	Code    int    `json:"code"`
}

func (r *RestError) Error() string {
	return r.Message
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

func NewInternalServerError(message string) *RestError {
	return &RestError{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestError {
	return &RestError{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *RestError {
	return &RestError{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}
