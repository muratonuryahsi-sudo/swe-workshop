package shared

import (
	"encoding/json"
	"errors"
	"net/http"
)

// errorResponse ist das JSON-Format fuer Fehlerantworten.
type errorResponse struct {
	Message string `json:"message"`
}

// WriteJSON schreibt body als JSON mit dem gegebenen Statuscode.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// WriteError mappt err auf den passenden HTTP-Statuscode und schreibt eine JSON-Fehlerantwort.
func WriteError(w http.ResponseWriter, err error) {
	var notFound *NotFoundError
	var validation *ValidationError
	var versionOutdated *VersionOutdatedError
	var unauthorized *UnauthorizedError
	var forbidden *ForbiddenError

	status := http.StatusInternalServerError
	switch {
	case errors.As(err, &notFound):
		status = http.StatusNotFound
	case errors.As(err, &validation):
		status = http.StatusUnprocessableEntity
	case errors.As(err, &versionOutdated):
		status = http.StatusConflict
	case errors.As(err, &unauthorized):
		status = http.StatusUnauthorized
	case errors.As(err, &forbidden):
		status = http.StatusForbidden
	}

	WriteJSON(w, status, errorResponse{Message: err.Error()})
}
