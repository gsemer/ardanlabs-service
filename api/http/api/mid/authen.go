package mid

import (
	"context"
	"net/http"

	"github.com/gsemer/ardanlabs-service/app/api/auth"
	"github.com/gsemer/ardanlabs-service/app/api/authclient"
	"github.com/gsemer/ardanlabs-service/app/api/mid"
	"github.com/gsemer/ardanlabs-service/foundation/logger"
	"github.com/gsemer/ardanlabs-service/foundation/web"
)

// AuthenticateService validates authentication via the auth service.
func Authenticate(log *logger.Logger, client *authclient.Client) web.MidHandler {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Authenticate(ctx, log, client, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}

// Bearer processes the authentication requirements locally.
func Bearer(auth *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Bearer(ctx, auth, r.Header.Get("authorization"), hdl)
		}

		return h
	}

	return m
}

// Basic processes the authentication requirements locally.
func Basic(auth *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Basic(ctx, hdl)
		}

		return h
	}

	return m
}
