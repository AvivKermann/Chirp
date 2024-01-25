package jsonResponse

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	Valid bool `json:"valid"`
}

func ResponedWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	errorStruct := ErrorResponse{
		Error: message,
	}

	data, err := json.Marshal(errorStruct)

	if err != nil {
		log.Printf("error marshalling JSON %s", err)
		statusCode = http.StatusInternalServerError
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}

func ResponedWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)

	if err != nil {
		ResponedWithError(w, http.StatusInternalServerError, "error marshalling json")
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}
