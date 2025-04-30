package httphelpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, status int, errorMsg string) {
	respBody := errorResponse{
		Error: errorMsg,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling error response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
