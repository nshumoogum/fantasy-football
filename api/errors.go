package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/nshumoogum/fantasy-football/models"

	errs "github.com/nshumoogum/fantasy-football/apierrors"
)

// ErrorResponse sets the structured error message in the http response body
func ErrorResponse(ctx context.Context, w http.ResponseWriter, status int, errorResponse *models.ErrorResponse) {
	b, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, errs.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		log.Event(ctx, "failed to write error response body", log.ERROR, log.Error(err), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
