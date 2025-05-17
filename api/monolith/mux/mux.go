package mux

import (
	"os"

	"github.com/andrew-hayworth22/critiquefy-service/api/monolith/route/sys"
	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

// WebAPI constructs a web app with all routes bound to it
func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown)

	sys.Routes(app)

	return app
}
