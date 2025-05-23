package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// Authorize is HTTP middleware that asserts that the user making the request has a role
func Authorize(auth *auth.Auth, rule string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Authorize(ctx, auth, rule, hdl)
		}

		return h
	}

	return m
}
