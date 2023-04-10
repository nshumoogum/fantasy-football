package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

// API represents an object containing configurable variables to connect to fantasy football API
type API struct {
	Client *http.Client
	URI    string
}

func (api *API) makeGetRequest(ctx context.Context, method, url string) (io.ReadCloser, error) {
	logData := log.Data{"url": url, "method": method}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Event(ctx, "failed to create request for fantasy football api", log.ERROR, log.Error(err), logData)
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.Client.Do(req)
	if err != nil {
		log.Event(ctx, "failed to action fantasy football api", log.ERROR, log.Error(err), logData)
		return nil, err
	}

	return resp.Body, nil
}
