package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// metrics defines all of the metrics we are tracking for debugging/profiling
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// Initializes metrics singleton
var m metrics

func init() {
	m = metrics{
		expvar.NewInt("goroutines"),
		expvar.NewInt("requests"),
		expvar.NewInt("errors"),
		expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

// Set updates the context with a pointer to the metrics data
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &m)
}

// UpdateGoroutines sets the goroutine value in the metrics data
func UpdateGoroutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		g := int64(runtime.NumGoroutine())
		v.goroutines.Set(g)
		return g
	}

	return 0
}

// AddRequest increments the requests value in the metrics data
func AddRequest(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.requests.Add(1)
		return v.requests.Value()
	}

	return 0
}

// AddError increments the errors value in the metrics data
func AddError(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.errors.Add(1)
		return v.errors.Value()
	}

	return 0
}

// AddPanic increments the panics value in the metrics data
func AddPanic(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.panics.Add(1)
		return v.panics.Value()
	}

	return 0
}
