package mid

import (
	"context"

	"github.com/andrew-hayworth22/critiquefy-service/app/metrics"
)

// Metrics is middleware that updates our metrics with error, request, and goroutine data
func Metrics(ctx context.Context, handler Handler) error {
	ctx = metrics.Set(ctx)

	err := handler(ctx)

	n := metrics.AddRequest(ctx)

	// Update goroutine count every 1000 requests
	if n%1000 == 0 {
		metrics.UpdateGoroutines(ctx)
	}

	if err != nil {
		metrics.AddError(ctx)
	}

	return err
}
