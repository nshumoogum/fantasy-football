package csv

import (
	"context"
	"os"

	"github.com/ONSdigital/log.go/log"
)

func writeToFile(ctx context.Context, connection *os.File, filename string, line string) error {
	logData := log.Data{
		"filename": filename,
		"line":     line,
	}
	_, err := connection.WriteString(line + "\n")
	if err != nil {
		log.Event(ctx, "error writing to file", log.ERROR, log.Error(err), logData)
		return err
	}

	return nil
}
