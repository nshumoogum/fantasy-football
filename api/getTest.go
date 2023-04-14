package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

type TestData struct {
	TestName string `json:"test_name"`
	TestAge  int    `json:"test_age"`
}

func (api *FantasyFootballAPI) getTest(w http.ResponseWriter, req *http.Request) {
	log.Event(req.Context(), "get test endpoint called", log.INFO)

	testData := &TestData{
		TestName: "John",
		TestAge:  24,
	}

	data, err := json.Marshal(testData)
	if err != nil {
		log.Event(req.Context(), "failed to marshal the response data due to", log.ERROR, log.Error(err))
		http.Error(w, "failed to encode data", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		log.Event(req.Context(), "writing response failed", log.ERROR, log.Error(err))
		http.Error(w, "Failed to write http response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

	return
}
