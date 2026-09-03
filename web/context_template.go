package web

// 模板渲染
import (
	"bytes"
	"io"
	"maps"

	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/web/template"
	"github.com/baiy/gbf/web/template/funcs"
)

// TemplateHandlerFunc 模板渲染请求处理方法
type TemplateHandlerFunc func(c *TemplateContext) error

type TemplateContext struct {
	*Context
	template     *template.Engine
	theme        string
	templateVars map[string]any

	contentType string
}

// Assign 模板变量注入
func (c *TemplateContext) Assign(name string, value any) *TemplateContext {
	c.templateVars[name] = value
	return c
}

// Assigns 模板变量注入
func (c *TemplateContext) Assigns(values map[string]any) *TemplateContext {
	maps.Copy(c.templateVars, values)
	return c
}

// Fetch 模板渲染
func (c *TemplateContext) Fetch(tpl string, themes ...string) error {
	if len(themes) == 0 {
		themes = append(themes, c.theme)
	}
	// 响应拦截 防止中途出现错误时, 提前输出响应
	wr := new(bytes.Buffer)
	err := c.template.Fetch(wr, themes[0], tpl, c.templateVars)
	if err != nil {
		return err
	}
	c.Response.ContentType(c.contentType)
	_, err = io.Copy(c.Response.Http(), wr)
	return err
}

// ContentType 设置响应内容类型
func (c *TemplateContext) ContentType(t string) *TemplateContext {
	c.contentType = t
	return c
}

// Exists 模板是否存在
func (c *TemplateContext) Exists(theme, tpl string) bool {
	return c.template.Exists(theme, tpl)
}

// Breadcrumb 面包屑导航
func (c *TemplateContext) Breadcrumb(title string, urls ...string) *TemplateContext {
	breadcrumbs := c.templateVars["breadcrumbs"].([]*breadcrumbItem)
	item := &breadcrumbItem{
		Title: title,
		Last:  true,
		First: len(breadcrumbs) == 1,
	}
	if len(urls) > 0 {
		item.Url = c.Url(urls[0])
	}
	for _, nav := range breadcrumbs {
		nav.Last = false
	}
	return c.Assign("breadcrumbs", append(breadcrumbs, item))
}

// Title 设置标题
func (c *TemplateContext) Title(v string, isBreadcrumbs ...bool) *TemplateContext {
	if len(isBreadcrumbs) > 0 && isBreadcrumbs[0] {
		c.Breadcrumb(v)
	}
	return c.Assign("title", v)
}

// Description 设置描述
func (c *TemplateContext) Description(v string) *TemplateContext {
	return c.Assign("description", v)
}

// Keywords 设置关键词
func (c *TemplateContext) Keywords(v string) *TemplateContext {
	return c.Assign("keywords", v)
}

// Success 成功跳转
// url 可为空 跳转地址，为空则跳转到回到上一页
// info 可选 提示信息 为空不展示
func (c *TemplateContext) Success(url, info string) error {
	return c.jump("success", info, url)
}

// Error 错误跳转
// url 可为空 跳转地址，为空则跳转到回到上一页
// info 可选 提示信息 为空不展示
func (c *TemplateContext) Error(url, info string) error {
	return c.jump("error", info, url)
}

// jump 结果页面
// t: success | error
func (c *TemplateContext) jump(t, info, url string) error {
	data := map[string]any{
		"info":    info,
		"url":     c.Url(url),
		"type":    t,
		"baseUrl": c.Web.option.baseUrl,
	}
	return Jump(c.Context, data)
}

func templateHandlerWrapper(
	w *Web,
	driver template.Driver,
	defaultTheme string,
	themeGetter func(*Context) string,
	templateFunc funcs.FuncMap,
) func(TemplateHandlerFunc) HandlerFunc {
	engine, err := template.NewEngine(driver, defaultTheme, w.Logger.With("template"), templateFunc)
	if err != nil {
		panic(err)
	}
	return func(h TemplateHandlerFunc) HandlerFunc {
		return func(c *Context) error {
			templateVars := map[string]any{
				"title":       "",
				"description": "",
				"keywords":    "",
				"baseUrl":     c.Web.option.baseUrl,
				"breadcrumbs": []*breadcrumbItem{},
			}
			maps.Copy(templateVars, c.Web.option.templateVars)

			context := &TemplateContext{Context: c, template: engine, templateVars: templateVars}
			if themeGetter != nil {
				context.theme = themeGetter(c)
			}

			err := errors.PanicToError(func() error {
				return h(context)
			})

			if err != nil {
				c.Logger.Errorf("template controller error:%v", err)
				// 业务异常处理
				if e := errors.IsBusiness(err); e != nil {
					return context.Error(e.Url, e.Message)
				}

				// 500 内部错误
				return Http500(c, err)
			}
			return nil
		}
	}
}

type breadcrumbItem struct {
	Title string
	Url   string
	Last  bool
	First bool
}
