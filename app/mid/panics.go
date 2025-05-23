package mid

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/andrew-hayworth22/critiquefy-service/app/metrics"
)

// Panics is middleware that catches and handles panics
func Panics(ctx context.Context, handler Handler) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

			metrics.AddPanic(ctx)
		}
	}()

	return handler(ctx)
}
