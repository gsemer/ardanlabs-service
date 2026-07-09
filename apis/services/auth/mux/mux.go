// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/gsemer/ardanlabs-service/apis/services/api/mid"
	"github.com/gsemer/ardanlabs-service/apis/services/auth/route/authapi"
	"github.com/gsemer/ardanlabs-service/apis/services/auth/route/checkapi"
	"github.com/gsemer/ardanlabs-service/business/api/auth"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
	"github.com/gsemer/ardanlabs-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(build string, log *logger.Logger, db *sqlx.DB, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(build, app, log, db)
	authapi.Routes(app, auth)

	return app
}
