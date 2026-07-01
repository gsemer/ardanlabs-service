// Package web contains a small framework extension.
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// A Handler is a type that handles a http request within our own little mini framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

// HandleFunc sets a handler function
func (a *App) HandleFunc(pattern string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		// PUT ANY CODE WE WANT HERE

		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
			return
		}

		// PUT ANY CODE WE WANT HERE

	}

	a.ServeMux.HandleFunc(pattern, h)
}
