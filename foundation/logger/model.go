package logger

import (
	"context"
	"log/slog"
	"time"
)

// Level represents a logging level
type Level slog.Level

// Set of defined levels
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Record represents data to be logged
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

// toRecord converts a slog record to our record
func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())

	f := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value.Any()
		return true
	}
	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: atts,
	}
}

// EventFn represents a function to be executed when configured against a log level
type EventFn func(ctx context.Context, r Record)

// Events represents an assignment of a function to each log level
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
