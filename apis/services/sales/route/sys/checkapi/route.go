package checkapi

import (
	"github.com/gsemer/ardanlabs-service/apis/services/api/mid"
	"github.com/gsemer/ardanlabs-service/app/api/authclient"
	"github.com/gsemer/ardanlabs-service/business/api/auth"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
	"github.com/gsemer/ardanlabs-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Routes adds specific routes for this group.
func Routes(build string, app *web.App, log *logger.Logger, db *sqlx.DB, authClient *authclient.Client) {
	authen := mid.AuthenticateService(log, authClient)
	authAdminOnly := mid.AuthorizeService(log, authClient, auth.RuleAdminOnly)

	api := newAPI(build, log, db)

	app.HandleFuncNoMiddleware("GET /liveness", api.liveness)
	app.HandleFuncNoMiddleware("GET /readiness", api.readiness)
	app.HandleFunc("GET /testerror", api.testError)
	app.HandleFunc("GET /testpanic", api.testPanic)
	app.HandleFunc("GET /testauth", api.liveness, authen, authAdminOnly)
}
