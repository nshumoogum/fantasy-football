package api

import (
	"context"
	"io"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// FantasyFootballAPI manages downloading of fpl data
type FantasyFootballAPI struct {
	Router *mux.Router
	FPLURL string
}

// NewFantasyFootballAPI create a new Fantasy Football API instance and register the API routes based on the application configuration.
func NewFantasyFootballAPI(ctx context.Context, fplURL string, router *mux.Router) *FantasyFootballAPI {
	api := &FantasyFootballAPI{
		Router: router,
		FPLURL: fplURL,
	}

	log.Event(ctx, "API and routing setup", log.INFO)

	api.Router.HandleFunc("/test", api.getTest).Methods("GET")
	api.Router.HandleFunc("/league-id/{id}/week/{event-week}", api.getDownload).Methods("GET")

	return api
}

// DrainBody drains the body of the given HTTP request
func DrainBody(r *http.Request) {

	if r.Body == nil {
		return
	}

	_, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		log.Event(r.Context(), "error draining request body", log.Error(err))
	}

	err = r.Body.Close()
	if err != nil {
		log.Event(r.Context(), "error closing request body", log.Error(err))
	}
}
