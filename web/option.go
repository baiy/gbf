package web

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/baiy/gbf/util"
	"github.com/baiy/gbf/web/template"
	"github.com/baiy/gbf/web/template/funcs"
	"github.com/gorilla/sessions"
)

type Option struct {
	// 系统入口目录
	baseUrl string
	// 服务端口
	port uint

	// Basic验证 username:password
	basicAuth string

	// 模板引擎驱动
	templateDriver func() template.Driver
	// 模板引擎默认主题
	templateDefaultTheme string
	// 模板引擎动态主题获取器
	templateThemeGetter func(*Context) string
	// 模板全局变量
	templateVars map[string]any
	// 自定义模板函数
	templateFunc funcs.FuncMap

	// info页面url
	infoUrl string

	// 是否启用session
	session bool
	// session 存储
	sessionStore func() sessions.Store

	// 是否启用
	accessLog bool
}

func WithPort(port uint) util.WithOption[Option] {
	return func(o *Option) {
		o.port = port
	}
}

func WithBasicAuth(username, password string) util.WithOption[Option] {
	return func(o *Option) {
		o.basicAuth = fmt.Sprintf("%s:%s", username, password)
	}
}

func WithTemplateDriver(driver func() template.Driver) util.WithOption[Option] {
	return func(o *Option) {
		o.templateDriver = driver
	}
}

func WithTemplateDefaultTheme(defaultTheme string) util.WithOption[Option] {
	return func(o *Option) {
		o.templateDefaultTheme = defaultTheme
	}
}

func WithTemplateThemeGetter(getter func(*Context) string) util.WithOption[Option] {
	return func(o *Option) {
		o.templateThemeGetter = getter
	}
}

func WithInfoUrl(url string) util.WithOption[Option] {
	return func(o *Option) {
		o.infoUrl = url
	}
}

func WithSession() util.WithOption[Option] {
	return func(o *Option) {
		o.session = true
	}
}

func WithSessionStore(store func() sessions.Store) util.WithOption[Option] {
	return func(o *Option) {
		o.sessionStore = store
	}
}

func WithTemplateVars(vars map[string]any) util.WithOption[Option] {
	return func(o *Option) {
		o.templateVars = vars
	}
}

func WithBaseUrl(path string) util.WithOption[Option] {
	return func(o *Option) {
		path := strings.TrimSuffix(strings.TrimPrefix(cmp.Or(path, "/"), "/"), "/")
		if path == "/" || path == "" {
			o.baseUrl = "/"
		}
		o.baseUrl = "/" + path + "/"
	}
}

func WithTemplateFunc(value funcs.FuncMap) util.WithOption[Option] {
	return func(o *Option) {
		o.templateFunc = value
	}
}

func WithAccessLog() util.WithOption[Option] {
	return func(o *Option) {
		o.accessLog = true
	}
}
