package api

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	errs "github.com/nshumoogum/fantasy-football/apierrors"
	"github.com/nshumoogum/fantasy-football/csv"
	"github.com/nshumoogum/fantasy-football/handlers"
	"github.com/nshumoogum/fantasy-football/models"
)

func (api *FantasyFootballAPI) getDownload(w http.ResponseWriter, req *http.Request) {
	log.Event(req.Context(), "get download endpoint called", log.INFO)

	defer DrainBody(req)
	ctx := req.Context()

	vars := mux.Vars(req)
	leagueID := vars["id"]
	eventWeek := vars["event-week"]
	logData := log.Data{"league_id": leagueID, "event_week": eventWeek}

	var errorObjects []*models.ErrorObject

	// Validate input
	week, errObjects := models.ValidateWeek(ctx, eventWeek, logData)
	if errObjects != nil {
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errObjects})
		return
	}

	filename := "league-id-" + leagueID + "-week-" + eventWeek + ".csv"

	file, err := os.Create(filename)
	if err != nil {
		log.Event(ctx, "error creating file", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
		return
	}
	defer file.Close()

	client := &http.Client{}
	fplAPI := &handlers.API{
		Client: client,
		URI:    api.FPLURL,
	}

	league, err := fplAPI.GetLeague(ctx, leagueID, 1)
	if err != nil {
		log.Event(ctx, "error getting league data", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
		return
	}

	if err = csv.CreateLeague(ctx, file, filename, league.League.Name, league.League.ID); err != nil {
		log.Event(ctx, "error adding league data", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
		return
	}

	if err = fplAPI.GetTeams(ctx, file, filename, week, league.Standings.Results); err != nil {
		log.Event(ctx, "error retrieving and adding team data", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
		return
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Event(ctx, "error retrieving file statistics", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
		return
	}

	openfile, err := os.Open(filename)
	if err != nil {
		log.Event(ctx, "failed to open file", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
	}

	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// Send the headers before sending the file
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", fileSize)

	// Send the file
	_, err = io.Copy(w, openfile)
	if err != nil {
		log.Event(ctx, "failed to download file", log.ERROR, log.Error(err), logData)
		errorObjects = append(errorObjects, &models.ErrorObject{Error: errs.ErrInternalServer.Error()})
		ErrorResponse(ctx, w, http.StatusInternalServerError, &models.ErrorResponse{Errors: errorObjects})
	}

	// TODO shouldn't need to remove files off disc (we shouldn't be storing files on disc in the first place)
	if err := os.Remove(filename); err != nil {
		log.Event(ctx, "error deleting file off disc", log.WARN, log.Error(err), logData)
	}

	log.Event(ctx, "download fpl csv: request successful", log.INFO, logData)
}
