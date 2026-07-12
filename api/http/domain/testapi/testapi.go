package testapi

import (
	"context"
	"math/rand/v2"
	"net/http"

	"github.com/gsemer/ardanlabs-service/app/api/errs"
	"github.com/gsemer/ardanlabs-service/foundation/web"
)

type api struct{}

func newAPI() *api {
	return &api{}
}

func (api *api) testError(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.IntN(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "this message is trusted")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, r, status, http.StatusOK)
}

func (api *api) testPanic(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.IntN(100); n%2 == 0 {
		panic("WE ARE PANICKING!!!")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, r, status, http.StatusOK)
}

func (api *api) testAuth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, r, status, http.StatusOK)
}
