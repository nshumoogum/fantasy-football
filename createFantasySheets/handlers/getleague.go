package handlers

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/ONSdigital/go-ns/log"
	"github.com/nshumoogum/fantasy-football/createFantasySheets/models"
	"github.com/pkg/errors"
)

// GetLeague retrieves league data from fantasy football site
func (api *API) GetLeague(ctx context.Context, leagueID string, page int) (*models.Resource, error) {
	method := "GET"
	pageString := strconv.Itoa(page)
	path := api.URI + "/leagues-classic-standings/" + leagueID + "?ls-page=" + pageString
	logData := log.Data{"url": path, "method": method}

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to create url for api call"), logData)
		return nil, err
	}
	path = URL.String()
	logData["url"] = path

	b, err := api.makeGetRequest(ctx, method, path)
	if err != nil {
		return nil, err
	}

	var resource models.Resource
	if err = json.Unmarshal(b, &resource); err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "unable to unmarshal bytes into league resource"), logData)
		return nil, err
	}

	// Check and retrieve other pages
	if resource.Standings != nil && resource.Standings.HasNext {
		page++
		next, err := api.GetLeague(ctx, leagueID, page)
		if err != nil {
			return nil, err
		}

		resource.Standings.Results = append(resource.Standings.Results, next.Standings.Results...)
	}

	log.Info("league data", log.Data{"resource": resource})

	return &resource, nil
}
