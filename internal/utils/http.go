package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func AnswerWithJSON(w http.ResponseWriter, payload any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AnswerWithInternalServerError(w http.ResponseWriter, err error) {
	AnswerWithJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
}

func AnswerWithNotFoundError(w http.ResponseWriter, err error) {
	AnswerWithJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
}

func AnswerWithBadRequestError(w http.ResponseWriter, err error) {
	AnswerWithJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, data any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type is not application/json")
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&data); err != nil {
		return err
	}

	return nil
}

func SetSSEResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
