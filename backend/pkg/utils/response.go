package utils

import (
    "encoding/json"
    "net/http"
    "time"
)

type Response struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Message   string      `json:"message,omitempty"`
    Error     *ErrorData  `json:"error,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}

type ErrorData struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    response := Response{
        Success:   true,
        Data:      data,
        Message:   message,
        Timestamp: time.Now(),
    }

    json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, code string, message string, details interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    response := Response{
        Success: false,
        Error: &ErrorData{
            Code:    code,
            Message: message,
            Details: details,
        },
        Timestamp: time.Now(),
    }

    json.NewEncoder(w).Encode(response)
}
