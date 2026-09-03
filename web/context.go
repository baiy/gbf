package web

// 通用 请求
import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/log"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

// HandlerFunc 请求处理方法
type HandlerFunc func(c *Context) error

type Context struct {
	echo   *echo.Context
	Logger log.Logger
	Web    *Web

	Request  *Request
	Response *Response

	session *Session

	start time.Time

	lock sync.RWMutex
}

// Ctx 获取请求的 context
func (c *Context) Ctx() context.Context {
	return c.Request.Http().Context()
}

func (c *Context) Get[T any](key string) (T, bool) {
	v, err := echo.ContextGet[T](c.echo, key)
	if err != nil {
		if errors.Is(err, echo.ErrNonExistentKey) {
			return v, false
		}
		panic(err)
	}
	return v, true
}

func (c *Context) Set[T any](key string, value T) {
	c.echo.Set(key, value)
}

func (c *Context) Session() *Session {
	if c.session != nil {
		return c.session
	}
	if !c.Web.option.session {
		panic("session is not enabled")
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	s, err := echo.ContextGet[sessions.Store](c.echo, "_session_store")
	if err != nil {
		panic(err)
	}
	session, err := s.Get(c.Request.Http(), strings.ToUpper(fmt.Sprintf("%s_session_id", constant.AppName)))
	if err != nil {
		c.Logger.Warnf("session get error: %v", err)
	}
	c.session = &Session{c: c, sess: session}
	return c.session
}

func (c *Context) Url(url string) string {
	if url == "" {
		return ""
	}
	if !strings.Contains(url, "://") {
		return c.Web.option.baseUrl + strings.TrimPrefix(url, "/")
	}
	return url
}

// Duration 获取请求处理耗时
func (c *Context) Duration() time.Duration {
	return time.Since(c.start)
}

// Start 获取请求开始时间
func (c *Context) Start() time.Time {
	return c.start
}

func (c *Context) dispatch(fn func() error) (err error) {
	defer func() {
		// 捕获 panic
		if r := recover(); r != nil {
			tmpErr, ok := r.(error)
			if ok {
				err = tmpErr
				return
			}
			panic(r)
		}
	}()
	return fn()
}

func newContext(c *echo.Context) *Context {
	v, err := echo.ContextGet[*Context](c, "c")
	if err != nil {
		if errors.Is(err, echo.ErrNonExistentKey) {
			w := c.Get("w").(*Web)
			ctx := &Context{echo: c, Web: w, Logger: w.Logger, start: time.Now()}
			ctx.Request = &Request{c: ctx, requestParam: map[string]any{}}
			ctx.Response = &Response{c: ctx}
			c.Set("c", ctx)
			return ctx
		}
		panic(err)
	}
	return v
}

func handlerWrapper(h HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		err := errors.PanicToError(func() error {
			return h(c)
		})

		if err != nil {
			c.Logger.Errorf("controller error:%v", err)
			// 业务异常处理
			if e := errors.IsBusiness(err); e != nil {
				return Http500(c, err)
			}
			return Http500(c, fmt.Errorf("系统错误"))
		}
		return nil
	}
}
