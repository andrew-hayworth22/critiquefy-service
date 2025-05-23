package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// Panics is HTTP middleware that catches and handles panics
func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Panics(ctx, hdl)
		}

		return h
	}

	return m
}
