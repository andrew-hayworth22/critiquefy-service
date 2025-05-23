package auth

import (
	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
)

func Routes(app *web.App, a *auth.Auth) {
	app.Handle("POST /auth/login", login)
}
