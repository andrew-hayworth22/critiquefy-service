package mid

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefy-service/app/errs"
	"github.com/andrew-hayworth22/critiquefy-service/app/mid"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// codeStatus maps application errors to HTTP error codes
var codeStatus [17]int

func init() {
	codeStatus[errs.OK.Value()] = http.StatusOK
	codeStatus[errs.Cancelled.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.NotFound.Value()] = http.StatusNotFound
	codeStatus[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[errs.Aborted.Value()] = http.StatusConflict
	codeStatus[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatus[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

// Errors is HTTP middleware that handles errors gracefully
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			if err := mid.Errors(ctx, log, hdl); err != nil {
				errs := err.(errs.Error)
				if err := web.Respond(ctx, w, errs, codeStatus[errs.Code.Value()]); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
