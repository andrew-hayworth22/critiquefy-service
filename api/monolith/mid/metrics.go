package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// Metrics is HTTP middleware that updates our metrics with error, request, and goroutine data
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Metrics(ctx, hdl)
		}

		return h
	}

	return m
}
