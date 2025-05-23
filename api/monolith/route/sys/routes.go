package sys

import (
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(app *web.App, build string, log *logger.Logger, db *pgxpool.Pool) {
	api := newAPI(build, log, db)

	app.HandleNoAppMiddleware("GET /liveness", api.liveness)
	app.HandleNoAppMiddleware("GET /readiness", api.readiness)
}
