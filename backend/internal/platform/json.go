package platform

import (
	"encoding/json"
	"net/http"
)

type ErrorEnvelope struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":{"code":"INTERNAL_ERROR","message":"Failed to encode JSON"}}`, http.StatusInternalServerError)
		}
	}
}

func WriteError(w http.ResponseWriter, status int, code, message string, details map[string]string) {
	env := ErrorEnvelope{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	WriteJSON(w, status, env)
}

func ReadJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
