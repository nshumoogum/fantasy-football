package csv

import (
	"context"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
)

func writeToFile(ctx context.Context, connection *os.File, filename string, line string) error {
	logData := log.Data{
		"filename": filename,
		"line":     line,
	}
	_, err := connection.WriteString(line + "\n")
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "error writing to file"), logData)
		return err
	}

	return nil
}
