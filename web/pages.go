package web

import (
	_ "embed"
	"html/template"

	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/errors"
)

//go:embed pages/500.tpl
var page500 string

//go:embed pages/404.tpl
var page404 string

//go:embed pages/jump.tpl
var pageJump string

var pagesTpl *template.Template

func pagesRender(ctx *Context, code int, name string, data any) error {
	ctx.Response.Status(code)
	return pagesTpl.ExecuteTemplate(ctx.Response.Http(), name, data)
}
func Http500(ctx *Context, err error) error {
	if env.IsProd() {
		err = nil
	}
	return pagesRender(ctx, 500, "500", map[string]any{"error": err})
}

func Http404(ctx *Context) error {
	return pagesRender(ctx, 404, "404", nil)
}

func Jump(ctx *Context, data map[string]any) error {
	return pagesRender(ctx, 200, "jump", data)
}

func init() {
	pagesTpl = template.New("system_pages")
	_, err1 := pagesTpl.New("500").Parse(page500)
	_, err2 := pagesTpl.New("404").Parse(page404)
	_, err3 := pagesTpl.New("jump").Parse(pageJump)
	err := errors.Join(err1, err2, err3)
	if err != nil {
		panic(err)
	}
}
