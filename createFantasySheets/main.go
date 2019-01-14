package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/nshumoogum/fantasy-football/createFantasySheets/csv"
	"github.com/nshumoogum/fantasy-football/createFantasySheets/handlers"
	"github.com/pkg/errors"
)

var (
	csvName string
	ctx     context.Context

	eventWeek = 0
	leagueID  = "66205"
	url       = "https://fantasy.premierleague.com/drf"
)

func main() {
	flag.StringVar(&csvName, "csv-name", csvName, "name of csv file")
	flag.IntVar(&eventWeek, "event-week", eventWeek, "the event week to create results")
	flag.StringVar(&leagueID, "league-id", leagueID, "id of league")
	flag.StringVar(&url, "url", url, "part of the url determined by the fantast football host and path")
	flag.Parse()

	logData := log.Data{
		"csv-name":   csvName,
		"event-week": eventWeek,
		"league-id":  leagueID,
		"url":        url,
	}

	ctx = context.Background()

	var missingFlags bool
	if eventWeek < 1 {
		log.ErrorCtx(ctx, errors.New("event-week is not set or is set to less than 1"), logData)
		missingFlags = true
	}

	if csvName == "" {
		log.ErrorCtx(ctx, errors.New("csv-name is not set"), logData)
		missingFlags = true
	}

	if missingFlags {
		os.Exit(1)
	}

	connection, err := os.OpenFile(csvName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "error opening file"), logData)
		os.Exit(1)
	}

	// https://fantasy.premierleague.com/drf/leagues-classic-standings/66205
	// https://fantasy.premierleague.com/drf/entry/564241/event/3/picks

	api := &handlers.API{
		Client: http.DefaultClient,
		URI:    url,
	}

	league, err := api.GetLeague(ctx, leagueID, 1)
	if err != nil {
		os.Exit(1)
	}

	if err = csv.CreateLeague(ctx, connection, csvName, league.League.Name, league.League.ID); err != nil {
		os.Exit(1)
	}

	if err = api.GetTeams(ctx, connection, csvName, eventWeek, league.Standings.Results); err != nil {
		os.Exit(1)
	}
}
