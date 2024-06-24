package http_helper

import (
	"encoding/json"
	"net/http"

	"github.com/danzBraham/cats-social/internal/helpers/validator"
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

func HandleErrorResponse(w http.ResponseWriter, status int, err error) {
	EncodeJSON(w, status, ResponseBody{
		Error:   http.StatusText(status),
		Message: err.Error(),
	})
}

func HandleSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	EncodeJSON(w, status, ResponseBody{
		Message: message,
		Data:    data,
	})
}

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	err := DecodeJSON(r, payload)
	if err != nil {
		HandleErrorResponse(w, http.StatusBadRequest, err)
		return err
	}

	err = validator.ValidatePayload(payload)
	if err != nil {
		EncodeJSON(w, http.StatusBadRequest, ResponseBody{
			Error:   "Request doesn't pass validation",
			Message: err.Error(),
		})
		return err
	}

	return nil
}
