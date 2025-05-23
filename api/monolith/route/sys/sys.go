package sys

import (
	"context"
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefy-service/business/data/sqldb"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

type api struct {
	build string
	log   *logger.Logger
	db    *pgxpool.Pool
}

func newAPI(build string, log *logger.Logger, db *pgxpool.Pool) *api {
	return &api{
		build: build,
		db:    db,
		log:   log,
	}
}

func (api *api) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

func (api *api) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK
	if err := sqldb.StatusCheck(ctx, api.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		api.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, data, statusCode)
}
