package mid

import (
	"context"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/app/errs"
)

// Authorize is middleware that asserts that the user making the request has a role
func Authorize(ctx context.Context, auth *auth.Auth, role string, handler Handler) error {
	claims := GetClaims(ctx)

	if err := auth.Authorize(ctx, claims, role); err != nil {
		return errs.Newf(errs.PermissionDenied, "unauthorized: claims [%v] role [%v]: %s", claims, role, err)
	}

	return handler(ctx)
}
