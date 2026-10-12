package web

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util"
	"github.com/baiy/gbf/util/slice"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CORSConfig = middleware.CORSConfig

// MiddlewareFunc 中间件函数
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func (m MiddlewareFunc) ToEcho() echo.MiddlewareFunc {
	return func(n echo.HandlerFunc) echo.HandlerFunc {
		next := func(c *Context) error {
			return n(c.echo)
		}
		return func(c *echo.Context) error {
			h := m(next)
			return h(newContext(c))
		}
	}
}

func middlewaresToEcho(middleware []MiddlewareFunc) []echo.MiddlewareFunc {
	return slice.Map(middleware, func(_ int, middleware MiddlewareFunc) echo.MiddlewareFunc {
		return middleware.ToEcho()
	})
}

// MiddlewareByEcho echo 中间件转换为 Web 中间件
func MiddlewareByEcho(middleware echo.MiddlewareFunc) MiddlewareFunc {
	return func(n HandlerFunc) HandlerFunc {
		next := func(c *echo.Context) error {
			return n(newContext(c))
		}
		return func(c *Context) error {
			h := middleware(next)
			return h(c.echo)
		}
	}
}
func EmptyMiddleware(next HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		return next(c)
	}
}

// BasicAuth 基础认证
// auth: username:password
func BasicAuth(auth string) echo.MiddlewareFunc {
	auths := strings.SplitN(auth, ":", 2)
	if len(auths) != 2 {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				return next(c)
			}
		}
	}
	return middleware.BasicAuth(func(c *echo.Context, username string, password string) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(auths[0])) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(auths[1])) == 1 {
			return true, nil
		}
		return false, nil
	})
}

// TrackId 跟踪id
func TrackId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// 重置跟踪id
			log.TrackReset(c.Response().Header().Get(echo.HeaderXRequestID))
			c.Response().Header().Set(echo.HeaderXRequestID, log.TrackId())
			return next(c)
		}
	}
}

func httpErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		if r, _ := echo.UnwrapResponse(c.Response()); r != nil && r.Committed {
			return
		}
		// 记录 panic 日志
		if _, ok := errors.AsType[*middleware.PanicStackError](err); ok {
			c.Logger().Error(fmt.Sprintf("http panic error %s:%v", c.Request().RequestURI, err))
		}

		code := http.StatusInternalServerError
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				code = tmp
			}
		}

		var cErr error
		if c.Request().Method == http.MethodHead { // Issue #608
			cErr = c.NoContent(code)
		} else {
			ctx := newContext(c)
			switch code {
			case 404:
				cErr = Http404(ctx)
			default:
				cErr = Http500(ctx, err)
			}
		}
		if cErr != nil {
			c.Logger().Error(fmt.Sprintf("web error handler failed to send error to client:%v", cErr))
		}
	}
}

// Cors 跨域中间件
func Cors(options ...util.WithOption[CorsOption]) MiddlewareFunc {
	option := &CorsOption{
		allowOrigins:     []string{"*"},
		allowCredentials: false,
		allowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	util.ApplyOption(option, options...)
	return MiddlewareByEcho(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     option.allowOrigins,
		AllowCredentials: option.allowCredentials,
		AllowMethods:     option.allowMethods,
		AllowHeaders: append([]string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAcceptEncoding,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
			echo.HeaderXRequestID,
		}, option.allowHeaders...),
		ExposeHeaders: append([]string{echo.HeaderXRequestID}, option.exposeHeaders...),
	}))
}

// 浏览器缓存中间件
func BrowserCache(days uint) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			c.Response.Header().Set(echo.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", 86400*days))
			return next(c)
		}
	}
}

type CorsOption struct {
	allowOrigins     []string
	allowCredentials bool
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
}

func WithCorsAllowOrigins(allowOrigins []string) util.WithOption[CorsOption] {
	return func(o *CorsOption) {
		o.allowOrigins = allowOrigins
	}
}

func WithCorsAllowCredentials() util.WithOption[CorsOption] {
	return func(o *CorsOption) {
		o.allowCredentials = true
	}
}

func WithCorsAllowMethods(allowMethods []string) util.WithOption[CorsOption] {
	return func(o *CorsOption) {
		o.allowMethods = allowMethods
	}
}

func WithCorsAllowHeaders(allowHeaders []string) util.WithOption[CorsOption] {
	return func(o *CorsOption) {
		o.allowHeaders = allowHeaders
	}
}

func WithCorsExposeHeaders(exposeHeaders []string) util.WithOption[CorsOption] {
	return func(o *CorsOption) {
		o.exposeHeaders = exposeHeaders
	}
}
