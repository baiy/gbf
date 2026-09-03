package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/manager"
	"github.com/baiy/gbf/util/cryptor"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func create(mo manager.Option[Option]) (*Web, error) {
	logger := log.With(fmt.Sprintf("web.%s", mo.Key))
	w := &Web{
		echo: echo.NewWithConfig(echo.Config{
			Logger:           logger.Slog(),
			HTTPErrorHandler: httpErrorHandler(),
			IPExtractor:      echo.ExtractIPFromXFFHeader(),
			JSONSerializer:   &JSONSerializer{},
		}),
		option: &mo.Value,
		Logger: logger,
	}
	w.echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("w", w)
			return next(c)
		}
	})

	// 错误处理
	w.echo.Use(middleware.Recover())
	// 跟踪id
	w.echo.Use(TrackId())
	// 请求日志
	w.echo.Use(middleware.RequestLogger())
	// session
	if w.option.session {
		w.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				c.Set("_session_store", w.option.sessionStore())
				return next(c)
			}
		})
	}
	// 基础认证
	if w.option.basicAuth != "" {
		w.echo.Use(BasicAuth(w.option.basicAuth))
	}

	// 通用路由
	w.Router = NewRouter(w, handlerWrapper)
	// api路由
	w.Api = NewRouter(w, ApiHandlerWrapper)
	// 服务间路由
	w.Service = NewRouter(w, ApiHandlerWrapper, func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			// 验证服务密钥
			if c.Request.Header().Get("X-Request-Service-Key") != cryptor.Sha256(constant.Secret) {
				c.Response.Status(http.StatusForbidden)
				return nil
			}
			return next(c)
		}
	})
	// 模板路由
	if w.option.templateDriver != nil {
		w.Template = NewRouter(w,
			templateHandlerWrapper(w,
				w.option.templateDriver(),
				w.option.templateDefaultTheme,
				w.option.templateThemeGetter,
				w.option.templateFunc,
			),
		)
	}
	closer.Register(fmt.Sprintf("web:%s", mo.Key), w)
	return w, nil
}

func start(w *Web) error {
	if w.option.infoUrl != "" {
		w.Router.GET(w.option.infoUrl, func(c *Context) error {
			return c.Response.JSON(map[string]any{
				"time":    time.Now().Format(time.DateTime + ".000000"),
				"gbf":     constant.BaseVersion,
				"app":     constant.AppName,
				"version": constant.AppVersion,
			})
		})
	}
	// 启动服务
	w.Start()
	return nil
}

var Manager = manager.New(
	"web",
	create,
	start,
	func() Option {
		return Option{
			baseUrl: "/",
			infoUrl: "info",
			port:    8080,
			sessionStore: func() sessions.Store {
				return sessions.NewCookieStore([]byte(constant.Secret))
			},

			templateVars: map[string]any{},
			templateFunc: map[string]any{},
		}
	},
)
