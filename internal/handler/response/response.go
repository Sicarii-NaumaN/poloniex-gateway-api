// Package helps to wrap response data and to set http status code

package response

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func OK(rw http.ResponseWriter, data interface{}) {
	result(rw, http.StatusOK, "Success", data)
}

func BadRequest(rw http.ResponseWriter, message string) {
	result(rw, http.StatusBadRequest, message, nil)
}

func Unauthenticated(rw http.ResponseWriter, message string) {
	result(rw, http.StatusUnauthorized, message, nil)
}

func NotFound(rw http.ResponseWriter, message string) {
	result(rw, http.StatusNotFound, message, nil)
}

func Forbidden(rw http.ResponseWriter, message string) {
	result(rw, http.StatusForbidden, message, nil)
}

func Internal(rw http.ResponseWriter, message string) {
	result(rw, http.StatusInternalServerError, message, nil)
}

func result(rw http.ResponseWriter, status int, message string, data any) {
	success := 200 <= status && status <= 300

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(Response[any]{
		Success: success,
		Message: message,
		Data:    data,
	})
}
