package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

func Authenticate(auth *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Authenticate(ctx, auth, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}
