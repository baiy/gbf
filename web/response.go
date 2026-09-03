package web

import (
	"cmp"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Response struct {
	c *Context

	wrapper *WrapperResponseWriter
}

func (r *Response) Http() http.ResponseWriter {
	return r.c.echo.Response()
}

func (r *Response) Header() http.Header {
	return r.Http().Header()
}

func (r *Response) XML(i any, codes ...int) error {
	return r.c.echo.XML(withDefaultOkStatus(codes), i)
}
func (r *Response) String(s string, codes ...int) error {
	return r.c.echo.String(withDefaultOkStatus(codes), s)
}
func (r *Response) Blob(contentType string, b []byte, codes ...int) error {
	return r.c.echo.Blob(withDefaultOkStatus(codes), contentType, b)
}
func (r *Response) JSON(i any, codes ...int) error {
	return r.c.echo.JSON(withDefaultOkStatus(codes), i)
}

func (r *Response) HTML(html string, codes ...int) error {
	return r.c.echo.HTML(withDefaultOkStatus(codes), html)
}

func (r *Response) NoContent(codes ...int) error {
	return r.c.echo.NoContent(withDefaultOkStatus(codes))
}

func (r *Response) ContentType(value string) {
	header := r.Http().Header()
	if header.Get(echo.HeaderContentType) == "" {
		header.Set(echo.HeaderContentType, cmp.Or(value, echo.MIMETextHTMLCharsetUTF8))
	}
}

func (r *Response) Status(code int) {
	r.Http().WriteHeader(code)
}

func (r *Response) Redirect(url string, codes ...int) error {
	url = r.c.Url(url)
	if len(codes) == 0 {
		return r.c.echo.Redirect(http.StatusFound, url)
	}
	return r.c.echo.Redirect(codes[0], url)
}

// Recorder 响应记录器
func (r *Response) Recorder() {
	_, ok := r.c.Get[*WrapperResponseWriter]("WrapperResponseWriter")
	if ok {
		return
	}
	r.wrapper = NewWrapperResponseWriter(r.c.echo.Response())
	r.c.echo.SetResponse(r.wrapper)
}

func (r *Response) Before(fn func()) {
	resp, err := echo.UnwrapResponse(r.c.echo.Response())
	if err != nil {
		r.c.Logger.Errorf("response before func error: %v", err)
		return
	}
	resp.Before(fn)
}

func (r *Response) After(fn func()) {
	resp, err := echo.UnwrapResponse(r.c.echo.Response())
	if err != nil {
		r.c.Logger.Errorf("response after func error: %v", err)
		return
	}
	resp.After(fn)
}

func (r *Response) GetResponseBody() string {
	return r.wrapper.String()
}

func withDefaultOkStatus(codes []int) int {
	if len(codes) == 0 {
		return http.StatusOK
	}
	return codes[0]
}
