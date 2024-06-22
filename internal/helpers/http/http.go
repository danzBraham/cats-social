package http_helper

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func DecodeJSON(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

func EncodeJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
