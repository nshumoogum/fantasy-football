package csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"os"

	"github.com/ONSdigital/log.go/log"
)

func writeFile(ctx context.Context, file *os.File, filename, line string) error {
	logData := log.Data{
		"filename": filename,
		"line":     line,
	}
	_, err := file.WriteString(line + "\n")
	if err != nil {
		log.Event(ctx, "error writing to file", log.ERROR, log.Error(err), logData)
		return err
	}

	return nil
}

// TODO use this method to create csv in memory instead of writing to disc and removing
func writeAll(records [][]string) ([]byte, error) {
	if records == nil || len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	err := csvWriter.WriteAll(records)
	if err != nil {
		return nil, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
