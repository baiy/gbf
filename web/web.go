package web

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/worker"
	"github.com/labstack/echo/v5"
)

// Web echo 框架重新封装
type Web struct {
	echo   *echo.Echo
	option *Option
	Logger log.Logger
	http   *http.Server

	Router   *Router[HandlerFunc]
	Api      *Router[ApiHandlerFunc]
	Service  *Router[ApiHandlerFunc]
	Template *Router[TemplateHandlerFunc]
}

func (w *Web) Use(middlewares ...MiddlewareFunc) {
	w.echo.Use(middlewaresToEcho(middlewares)...)
}

func (w *Web) Pre(middlewares ...MiddlewareFunc) {
	w.echo.Pre(middlewaresToEcho(middlewares)...)
}

func (w *Web) Start() {
	w.http = &http.Server{Addr: fmt.Sprintf(":%d", cmp.Or(w.option.port, 8080)), Handler: w.echo}
	worker.Async(func() {
		w.Logger.Infof("starting server on %s", w.http.Addr)
		if err := w.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			w.Logger.Errorf("failed to start server error: %s", err)
		}
	})
}

func (w *Web) Close() error {
	if w.http == nil {
		return nil
	}
	w.http.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return w.http.Shutdown(ctx)
}
