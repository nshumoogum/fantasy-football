package csv

import (
	"context"
	"os"
	"strconv"
)

// CreateLeague ...
func CreateLeague(ctx context.Context, connection *os.File, filename, name string, id int) (err error) {
	titleLine := name + " (" + strconv.Itoa(id) + ")\n"
	if err = writeToFile(ctx, connection, filename, titleLine); err != nil {
		return
	}

	headerLine := "Entry ID,Player,Team Name,Number of Transfers,Transfer Cost,Points"
	if err = writeToFile(ctx, connection, filename, headerLine); err != nil {
		return
	}

	return
}
