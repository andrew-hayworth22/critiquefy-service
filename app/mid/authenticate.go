package mid

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/app/errs"
	"github.com/google/uuid"
)

// Authenticate is middleware that processes an auth token and stores user data in the context
func Authenticate(ctx context.Context, auth *auth.Auth, authorization string, handler Handler) error {
	var err error
	parts := strings.Split(authorization, " ")

	if len(parts) != 2 {
		return errs.New(errs.Unauthenticated, errors.New("malformed token"))
	}

	switch parts[0] {
	case "Bearer":
		ctx, err = processJWT(ctx, auth, parts[1])
	}

	if err != nil {
		return err
	}

	return handler(ctx)
}

// processJWT processes information from a JWT Bearer token
func processJWT(ctx context.Context, auth *auth.Auth, token string) (context.Context, error) {
	claims, err := auth.Authenticate(ctx, token)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return ctx, errs.New(errs.Unauthenticated, errors.New("no subject claim"))
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = setUserId(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}
