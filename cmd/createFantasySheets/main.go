package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/ONSdigital/log.go/log"
	"github.com/nshumoogum/fantasy-football/csv"
	"github.com/nshumoogum/fantasy-football/handlers"
	"github.com/pkg/errors"
)

var (
	ctx context.Context

	defaultLeagueID      = "414763"
	defaultURL           = "https://fantasy.premierleague.com/api"
	defaultFileExtension = ".csv"

	fileExtension, leagueID, url string
	eventWeek                    int
)

func main() {
	flag.StringVar(&fileExtension, "file-extension", fileExtension, "file extension")
	flag.IntVar(&eventWeek, "event-week", eventWeek, "the event week to create results")
	flag.StringVar(&leagueID, "league-id", leagueID, "id of league")
	flag.StringVar(&url, "url", url, "part of the url determined by the fantast football host and path")
	flag.Parse()

	logData := log.Data{
		"event-week": eventWeek,
		"league-id":  leagueID,
		"url":        url,
	}

	ctx = context.Background()

	if eventWeek < 1 {
		log.Event(ctx, "event-week is not set or is set to less than 1", log.ERROR, log.Error(errors.New("flag event-week not set")), logData)
		os.Exit(1)
	}

	if fileExtension == "" {
		log.Event(ctx, "file-extension is not set, using default", log.WARN, log.Error(errors.New("flag file-extension not set")), logData)
		fileExtension = defaultFileExtension
	}

	if leagueID == "" {
		log.Event(ctx, "league-id is not set, using default", log.WARN, log.Error(errors.New("flag league-id not set")), logData)
		leagueID = defaultLeagueID
	}

	if url == "" {
		log.Event(ctx, "(fantasy football) url is not set, using default", log.WARN, log.Error(errors.New("flag url not set")), logData)
		url = defaultURL
	}

	filename := "files/week-" + strconv.Itoa(eventWeek) + fileExtension

	connection, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Event(ctx, "error opening file", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	// https://fantasy.premierleague.com/drf/leagues-classic-standings/66205
	// https://fantasy.premierleague.com/drf/entry/564241/event/3/picks
	client := &http.Client{}
	api := &handlers.API{
		Client: client,
		URI:    url,
	}

	league, err := api.GetLeague(ctx, leagueID, 1)
	if err != nil {
		os.Exit(1)
	}

	if err = csv.CreateLeague(ctx, connection, filename, league.League.Name, league.League.ID); err != nil {
		os.Exit(1)
	}

	if err = api.GetTeams(ctx, connection, filename, eventWeek, league.Standings.Results); err != nil {
		os.Exit(1)
	}
}
