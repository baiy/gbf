package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/baiy/gbf/env"
	"github.com/labstack/echo/v5"
)

type RouteInfo = echo.RouteInfo

func handlerWrap[H any](r *Router[H], h H) echo.HandlerFunc {
	return func(c *echo.Context) error {
		handler := r.handlerWrap(h)
		err := handler(newContext(c))
		if err != nil {
			r.web.Logger.Errorf("http request error:%s", err)
		}
		return err
	}
}

type Router[H any] struct {
	web         *Web
	handlerWrap func(h H) HandlerFunc

	middlewares []MiddlewareFunc
}

func (r *Router[H]) Add(method, path string, h H, middleware ...MiddlewareFunc) RouteInfo {
	return r.web.echo.Add(method, path, handlerWrap(r, h), middlewaresToEcho(r.injectGlobalMiddleware(middleware))...)
}

func (r *Router[H]) GET(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodGet, path, h, middlewares...)
}

func (r *Router[H]) POST(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodPost, path, h, middlewares...)
}

func (r *Router[H]) PUT(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodPut, path, h, middlewares...)
}

func (r *Router[H]) DELETE(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodDelete, path, h, middlewares...)
}

func (r *Router[H]) CONNECT(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodConnect, path, h, middlewares...)
}

func (r *Router[H]) HEAD(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodHead, path, h, middlewares...)
}

func (r *Router[H]) OPTIONS(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodOptions, path, h, middlewares...)
}

func (r *Router[H]) PATCH(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodPatch, path, h, middlewares...)
}

func (r *Router[H]) TRACE(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(http.MethodTrace, path, h, middlewares...)
}

func (r *Router[H]) RouteNotFound(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(echo.RouteNotFound, path, h, middlewares...)
}

func (r *Router[H]) Any(path string, h H, middlewares ...MiddlewareFunc) RouteInfo {
	return r.Add(echo.RouteAny, path, h, middlewares...)
}

func (r *Router[H]) Match(methods []string, path string, h H, middleware ...MiddlewareFunc) echo.Routes {
	return r.web.echo.Match(methods, path, handlerWrap(r, h), middlewaresToEcho(r.injectGlobalMiddleware(middleware))...)
}

func (r *Router[H]) Group(prefix string, middleware ...MiddlewareFunc) *Group[H] {
	return newGroup(r.web.echo.Group(prefix, middlewaresToEcho(r.injectGlobalMiddleware(middleware))...), r)
}

func (r *Router[H]) Static(pathPrefix, fsRoot string, middleware ...MiddlewareFunc) RouteInfo {
	return r.web.echo.Static(pathPrefix, fsRoot, r.staticMiddleware(middleware)...)
}

func (r *Router[H]) StaticFS(pathPrefix string, filesystem fs.FS, middleware ...MiddlewareFunc) RouteInfo {
	return r.web.echo.StaticFS(pathPrefix, filesystem, r.staticMiddleware(middleware)...)
}

// StaticEmbed 静态资源服务
func (r *Router[H]) StaticEmbed(
	pathPrefix string,
	efs embed.FS,
	efsPath string,
	devAbsPath string,
	middleware ...MiddlewareFunc,
) RouteInfo {
	// 本地环境可使用本地目录(绝对路径) 方便动态修改文件调试
	if devAbsPath != "" && env.IsLocal() {
		if filepath.IsAbs(devAbsPath) {
			info, err := os.Stat(devAbsPath)
			if err == nil && info.IsDir() {
				return r.Static(pathPrefix, devAbsPath, middleware...)
			}
			r.web.Logger.Errorf("StaticEmbed path error: %s", err)
		}
	}
	return r.StaticFS(pathPrefix, echo.MustSubFS(efs, efsPath), middleware...)
}

func (r *Router[H]) File(path, file string, middleware ...MiddlewareFunc) RouteInfo {
	return r.web.echo.File(path, file, r.staticMiddleware(middleware)...)
}

func (r *Router[H]) FileFS(path, file string, filesystem fs.FS, middleware ...MiddlewareFunc) RouteInfo {
	return r.web.echo.FileFS(path, file, filesystem, r.staticMiddleware(middleware)...)
}

// FileEmbed 静态资源服务(单一文件)
func (r *Router[H]) FileEmbed(path, file string, efs embed.FS, efsPath string, devAbsPath string, middleware ...MiddlewareFunc) RouteInfo {
	// 本地环境可使用本地目录(绝对路径) 方便动态修改文件调试
	if devAbsPath != "" && env.IsLocal() {
		if filepath.IsAbs(devAbsPath) {
			info, err := os.Stat(devAbsPath)
			if err == nil && info.IsDir() {
				return r.File(path, devAbsPath+"/"+file, middleware...)
			}
			r.web.Logger.Errorf("FileEmbed path error: %s", err)
		}
	}
	return r.FileFS(path, file, echo.MustSubFS(efs, efsPath), middleware...)
}

// 静态资源浏览器缓存
func staticBrowserCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", 86400*31))
		return next(c)
	}
}

// 静态文件中间件
func (r *Router[H]) staticMiddleware(middlewares []MiddlewareFunc) []echo.MiddlewareFunc {
	if env.IsProd() {
		// 生产环境启用静态文件浏览器缓存
		return append(middlewaresToEcho(r.injectGlobalMiddleware(middlewares)), staticBrowserCache)
	}
	return middlewaresToEcho(r.injectGlobalMiddleware(middlewares))
}

func (r *Router[H]) injectGlobalMiddleware(middlewares []MiddlewareFunc) []MiddlewareFunc {
	if len(r.middlewares) == 0 {
		return middlewares
	}
	return append(middlewares, r.middlewares...)
}

func NewRouter[H any](w *Web, h func(handler H) HandlerFunc, middlewares ...MiddlewareFunc) *Router[H] {
	return &Router[H]{web: w, handlerWrap: h, middlewares: middlewares}
}

type Group[H any] struct {
	eg *echo.Group
	r  *Router[H]
}

func (g *Group[H]) Use(middleware ...MiddlewareFunc) {
	g.eg.Use(middlewaresToEcho(middleware)...)
}

func (g *Group[H]) CONNECT(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodConnect, path, h, m...)
}

func (g *Group[H]) DELETE(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodDelete, path, h, m...)
}

func (g *Group[H]) GET(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodGet, path, h, m...)
}

func (g *Group[H]) HEAD(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodHead, path, h, m...)
}

func (g *Group[H]) OPTIONS(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodOptions, path, h, m...)
}

func (g *Group[H]) PATCH(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodPatch, path, h, m...)
}

func (g *Group[H]) POST(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodPost, path, h, m...)
}

func (g *Group[H]) PUT(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodPut, path, h, m...)
}

func (g *Group[H]) TRACE(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(http.MethodTrace, path, h, m...)
}

func (g *Group[H]) Any(path string, handler H, middleware ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(echo.RouteAny, path, handler, middleware...)
}

func (g *Group[H]) Match(methods []string, path string, h H, middleware ...MiddlewareFunc) echo.Routes {
	return g.eg.Match(methods, path, handlerWrap(g.r, h), middlewaresToEcho(middleware)...)
}

func (g *Group[H]) RouteNotFound(path string, h H, m ...MiddlewareFunc) echo.RouteInfo {
	return g.Add(echo.RouteNotFound, path, h, m...)
}

func (g *Group[H]) Add(method, path string, h H, middleware ...MiddlewareFunc) echo.RouteInfo {
	return g.eg.Add(method, path, handlerWrap(g.r, h), middlewaresToEcho(middleware)...)
}

func (g *Group[H]) Group(prefix string, middleware ...MiddlewareFunc) *Group[H] {
	return newGroup(g.eg.Group(prefix, middlewaresToEcho(middleware)...), g.r)
}

func newGroup[H any](eg *echo.Group, r *Router[H]) *Group[H] {
	return &Group[H]{eg: eg, r: r}
}
