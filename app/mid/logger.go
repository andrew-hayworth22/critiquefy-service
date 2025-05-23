package mid

import (
	"context"
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// Logger is middleware that logs before and after the request is processed
func Logger(ctx context.Context, log *logger.Logger, path string, rawQuery string, method string, remoteAddr string, handler Handler) error {
	v := web.GetValues(ctx)

	if rawQuery != "" {
		path = fmt.Sprintf("%s?%s", path, rawQuery)
	}

	log.Info(ctx, "request started", "method", method, "path", path, "remoteAddr", remoteAddr)

	err := handler(ctx)

	log.Info(ctx, "request completed", "method", method, "path", path, "remoteAddr", remoteAddr,
		"statuscode", v.StatusCode, "since", time.Since(v.Now).String())

	return err
}
