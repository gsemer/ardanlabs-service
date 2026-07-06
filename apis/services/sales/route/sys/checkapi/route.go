package checkapi

import (
	"github.com/gsemer/ardanlabs-service/apis/services/api/mid"
	"github.com/gsemer/ardanlabs-service/app/api/authclient"
	"github.com/gsemer/ardanlabs-service/business/api/auth"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
	"github.com/gsemer/ardanlabs-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, authClient *authclient.Client) {
	authen := mid.AuthenticateService(log, authClient)
	authAdminOnly := mid.AuthorizeService(log, authClient, auth.RuleAdminOnly)

	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
	app.HandleFunc("GET /testerror", testError)
	app.HandleFunc("GET /testpanic", testPanic)
	app.HandleFunc("GET /testauth", liveness, authen, authAdminOnly)
}
