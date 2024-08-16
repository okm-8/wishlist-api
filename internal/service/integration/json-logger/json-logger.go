package jsonLogger

import (
	"api/internal/model/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
)

type jsonRecord struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	LevelCode log.Level      `json:"levelCode"`
	Message   string         `json:"message"`
	Labels    map[string]any `json:"labels"`
}

func PrintLogRecord(ctx Context, record log.Record) error {
	labels := make(map[string]any)

	for _, label := range record.Labels() {
		labels[label.Key()] = label.Value()
	}

	_jsonRecord := jsonRecord{
		Timestamp: record.Timestamp().Format("2006-01-02 15:04:05.000"),
		Level:     record.Level().String(),
		LevelCode: record.Level(),
		Message:   record.Message(),
		Labels:    labels,
	}

	message, err := json.Marshal(_jsonRecord)

	if err != nil {
		_jsonRecord.Labels = map[string]any{
			"_":     "Error marshalling log record",
			"error": err.Error(),
		}

		message, _ = json.Marshal(_jsonRecord)
	}

	if ctx.ColoredLog() {
		var logRecordColor pterm.Color
		switch record.Level() {
		case log.Debug:
			logRecordColor = pterm.FgGray
		case log.Info:
			logRecordColor = pterm.FgCyan
		case log.Warning:
			logRecordColor = pterm.FgYellow
		case log.Error:
			logRecordColor = pterm.FgRed
		default:
			logRecordColor = pterm.FgDefault
		}

		message = []byte(logRecordColor.Sprint(string(message)))
	}

	_, err = ctx.LogOutput().Write(append(message, '\n'))

	if err != nil {
		return errors.New(fmt.Sprintf("Error writing log record: %s", err))
	}

	return nil
}
