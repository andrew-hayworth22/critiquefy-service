package web

import (
	"context"
	"time"
)

type ctxKey int

const key ctxKey = 1

// Values represents information stored in the context of each web request
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues retrieves all information from the context
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

// GetTraceID retrieves the TraceID from the context
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

// GetTime retrieves the Now timestamp from the context
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return time.Now()
	}
	return v.Now
}

// setStatusCode sets the StatusCode of the context
func setStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}
	v.StatusCode = statusCode
}

// setValues sets the Values of the context
func setValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}
