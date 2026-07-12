// Package all binds all the routes into the specified app.
package all

import (
	"github.com/gsemer/ardanlabs-service/api/http/api/mux"
	"github.com/gsemer/ardanlabs-service/api/http/domain/checkapi"
	"github.com/gsemer/ardanlabs-service/foundation/web"
)

func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})
}
