package handlers

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
)

// API represents an object containing configurable variables to connect to fantasy football API
type API struct {
	Client *http.Client
	URI    string
}

func (api *API) makeGetRequest(ctx context.Context, method, url string) ([]byte, error) {
	logData := log.Data{"url": url, "method": method}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to create request for fantasy football api"), logData)
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.Client.Do(req)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "Failed to action fantasy football api"), logData)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to read body from fantasy football api"), logData)
		return nil, err
	}

	return responseBody, nil
}
