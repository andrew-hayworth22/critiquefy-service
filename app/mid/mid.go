package mid

import (
	"context"
	"errors"

	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/google/uuid"
)

// Handler represents application logic
type Handler func(context.Context) error

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
)

// setClaims sets the claims in the context for later use
func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims fetches user claims from the context
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

// setUserId sets the user ID in the context for later use
func setUserId(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// GetUserId fetches the user ID from the context
func GetUserId(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}
	return v, nil
}
