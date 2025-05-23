package mux

import (
	"os"

	"github.com/andrew-hayworth22/critiquefy-service/api/monolith/mid"
	authAPI "github.com/andrew-hayworth22/critiquefy-service/api/monolith/route/auth"
	"github.com/andrew-hayworth22/critiquefy-service/api/monolith/route/sys"
	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// WebAPI constructs a web app with all routes bound to it
func WebAPI(log *logger.Logger, a *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	sys.Routes(app)
	authAPI.Routes(app, a)

	return app
}
