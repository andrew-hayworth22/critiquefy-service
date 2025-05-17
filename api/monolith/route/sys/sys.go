package sys

import "github.com/andrew-hayworth22/critiquefy-service/foundation/web"

func Routes(app *web.App) {
	app.HandleNoAppMiddleware("GET /liveness", liveness)
}
