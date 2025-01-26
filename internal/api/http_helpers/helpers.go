package http_helpers

import (
	"encoding/json"
	"net/http"

	"applicationDesignTest/pkg/log"
)

type HttpStatus string
type ErrorType string

var (
	StatusSuccess HttpStatus = "success"
	StatusError   HttpStatus = "error"

	ErrorTypeValidationError ErrorType = "validation error"
	ErrorTypeInternalError   ErrorType = "internal server error"
)

type SuccessResponse struct {
	Status HttpStatus `json:"status"`
	Data   any        `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  HttpStatus `json:"status"`
	Error   ErrorType  `json:"error"`
	Message string     `json:"message"`
}

func SendError(w http.ResponseWriter, statusCode int, errMsg string, errorType ErrorType) {
	w.Header().Set("Content-Type", "application/json")

	resp := ErrorResponse{
		Status:  StatusError,
		Error:   errorType,
		Message: errMsg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode error response", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

func SendSuccess(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := SuccessResponse{
		Status: StatusSuccess,
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode success response", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
