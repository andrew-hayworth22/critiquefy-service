package mid

import (
	"context"

	"github.com/andrew-hayworth22/critiquefy-service/app/errs"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
)

// Errors is middleware that handles errors gracefully
func Errors(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)
	if err == nil {
		return nil
	}

	log.Error(ctx, "message", "ERROR", err.Error())

	if errs.IsError(err) {
		return errs.GetError(err)
	}

	return errs.New(errs.Unknown, err)
}
