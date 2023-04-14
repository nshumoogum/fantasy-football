package api

import (
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

func (api *FantasyFootballAPI) getHealth(w http.ResponseWriter, req *http.Request) {
	log.Event(req.Context(), "get health endpoint called", log.INFO)

	w.WriteHeader(204)

	return
}
