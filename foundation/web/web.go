package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// Handler represents logic that can handle an HTTP request
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App represents an entrypoint to the application and configures ctx object
type App struct {
	*http.ServeMux
	shutdown      chan os.Signal
	appMiddleware []Middleware
}

// NewApp creates a new web application
func NewApp(shutdown chan os.Signal, appMiddleware ...Middleware) *App {
	return &App{
		ServeMux:      http.NewServeMux(),
		shutdown:      shutdown,
		appMiddleware: appMiddleware,
	}
}

// SignalShutdown shuts down the web app
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for an HTTP method and path and includes app middleware
func (a *App) Handle(pattern string, handler Handler, middleware ...Middleware) {
	handler = wrapMiddleware(middleware, handler)
	handler = wrapMiddleware(a.appMiddleware, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	a.HandleFunc(pattern, h)
}

// HandleNoAppMiddleware sets a handler function for an HTTP method and path and excludes app middleware
func (a *App) HandleNoAppMiddleware(pattern string, handler Handler, middleware ...Middleware) {
	handler = wrapMiddleware(middleware, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	a.HandleFunc(pattern, h)
}

// validateError checks if the error requires a system shutdown
func validateError(err error) bool {
	switch {
	case errors.Is(err, syscall.EPIPE):
		return false
	case errors.Is(err, syscall.ECONNRESET):
		return false
	}

	return true
}
