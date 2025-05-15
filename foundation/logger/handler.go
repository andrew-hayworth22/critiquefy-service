package logger

import (
	"context"
	"log/slog"
)

// logHandler wraps slog handler to capture the log level for event handling
type logHandler struct {
	handler slog.Handler
	events  Events
}

// newLogHandler creates a new logHandler
func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handles records at given level
// Handler ignores records of lower level
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns JSONHandler with h's attrs followed by attrs
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs), events: h.events}
}

// WithGroup returns new handler with group appended to receiver's groups
// Keys of subsequent attributes should be qualified by the sequence of group names
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), events: h.events}
}

// Handle calls respective handlers depending on Events config
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if h.events.Debug != nil {
			h.events.Debug(ctx, toRecord(r))
		}
	case slog.LevelError:
		if h.events.Error != nil {
			h.events.Error(ctx, toRecord(r))
		}
	case slog.LevelWarn:
		if h.events.Warn != nil {
			h.events.Warn(ctx, toRecord(r))
		}
	case slog.LevelInfo:
		if h.events.Info != nil {
			h.events.Info(ctx, toRecord(r))
		}
	}

	return h.handler.Handle(ctx, r)
}
