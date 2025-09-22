package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Status:  status,
		Message: msg,
	})
}

type IDParsingError struct {
	value string
}

func NewIDParsingError(value string) error {
	return &IDParsingError{value: value}
}

func (e *IDParsingError) Error() string {
	return fmt.Sprintf("parsing id: invalid value '%v'", e.value)
}
