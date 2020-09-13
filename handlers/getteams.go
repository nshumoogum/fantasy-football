package handlers

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/ONSdigital/log.go/log"
	"github.com/nshumoogum/fantasy-football/csv"
	"github.com/nshumoogum/fantasy-football/models"
)

// GetTeams retrieves a list of team data and adds the data to csv
func (api *API) GetTeams(ctx context.Context, connection *os.File, filename string, event int, teams []*models.Result) error {
	method := "GET"

	// Loop through list of teams
	for _, team := range teams {
		path := api.URI + "/entry/" + strconv.Itoa(team.Entry) + "/event/" + strconv.Itoa(event) + "/picks/"
		logData := log.Data{"url": path, "method": method}

		URL, err := url.Parse(path)
		if err != nil {
			log.Event(ctx, "failed to create url for api call", log.ERROR, log.Error(err), logData)
			return err
		}
		path = URL.String()
		logData["url"] = path

		b, err := api.makeGetRequest(ctx, method, path)
		if err != nil {
			return err
		}

		var teamData models.Team
		if err = json.Unmarshal(b, &teamData); err != nil {
			log.Event(ctx, "unable to unmarshal bytes into team data resource", log.ERROR, log.Error(err), logData)
			return err
		}

		if err = csv.CreateTeam(ctx, connection, filename, team, teamData.EntryHistory); err != nil {
			return err
		}
	}

	return nil
}
