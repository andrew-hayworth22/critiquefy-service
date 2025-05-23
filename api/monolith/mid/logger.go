package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// Logger is middleware that logs before and after the request is processed
func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Logger(ctx, log, r.URL.Path, r.URL.RawQuery, r.Method, r.RemoteAddr, hdl)
		}

		return h
	}

	return m
}
