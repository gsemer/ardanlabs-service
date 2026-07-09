// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/gsemer/ardanlabs-service/apis/services/api/mid"
	"github.com/gsemer/ardanlabs-service/apis/services/sales/route/sys/checkapi"
	"github.com/gsemer/ardanlabs-service/app/api/authclient"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
	"github.com/gsemer/ardanlabs-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(build string, log *logger.Logger, db *sqlx.DB, authClient *authclient.Client, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(build, mux, log, db, authClient)

	return mux
}
