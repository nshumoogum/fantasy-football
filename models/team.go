package models

import (
	"context"
	"strconv"

	"github.com/ONSdigital/log.go/log"
	errs "github.com/nshumoogum/fantasy-football/apierrors"
)

// Team ...
type Team struct {
	EntryHistory *EntryHistory `json:"entry_history"`
}

// EntryHistory ...
type EntryHistory struct {
	Points             int `json:"points"`
	EventTransfers     int `json:"event_transfers"`
	EventTransfersCost int `json:"event_transfers_cost"`
}

// ValidateWeek ...
func ValidateWeek(ctx context.Context, week string, logData log.Data) (w int, errorObjects []*ErrorObject) {
	w, err := strconv.Atoi(week)
	if err != nil {
		log.Event(ctx, "failed to convert week type from string to integer", log.ERROR, log.Error(err), logData)
		errorValues := map[string](string){"event-week": week}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrEventWeekValue.Error(), ErrorValues: errorValues})

		return
	}

	if w < 1 {
		log.Event(ctx, "failed validation: week value below 1", log.ERROR, log.Error(errs.ErrEventWeekValue), logData)
		errorValues := map[string](string){"event-week": week}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrEventWeekValue.Error(), ErrorValues: errorValues})
	}

	return
}
