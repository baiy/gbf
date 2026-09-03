package web

import (
	"bytes"
	"cmp"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/baiy/gbf/errors"
	"github.com/labstack/echo/v5"
)

type PathValues = echo.PathValues

type Request struct {
	c *Context

	requestParam map[string]any
}

// Http 获取原始请求
func (r *Request) Http() *http.Request {
	return r.c.echo.Request()
}

func (r *Request) Header() http.Header {
	return r.Http().Header
}

// Input 根据类型优先级获取请求参数
//
// Query参数(?id=) > 路径参数(/:id) > 表单参数
//
// 获取指定类型参数请求方法: Query Path Form
func (r *Request) Input[T comparable](name string, defaults ...T) T {
	return param(r, "", name, defaults...)
}

// Query 查询参数
func (r *Request) Query[T comparable](name string, defaults ...T) T {
	return param(r, "query", name, defaults...)
}

// Form 表单参数
func (r *Request) Form[T comparable](name string, defaults ...T) T {
	return param(r, "form", name, defaults...)
}

// Path 路径参数
func (r *Request) Path[T comparable](name string, defaults ...T) T {
	return param(r, "path", name, defaults...)
}

// Form 表单参数
func (r *Request) Forms() url.Values {
	v, err := r.c.echo.FormValues()
	if err != nil {
		panic(err)
	}
	return v
}

// Querys 查询参数
func (r *Request) Querys() url.Values {
	return r.c.echo.QueryParams()
}

// Paths 路径参数
func (r *Request) Paths() PathValues {
	return r.c.echo.PathValues()
}

// File 单文件上传
func (r *Request) File(name string) *multipart.FileHeader {
	file, err := r.c.echo.FormFile(name)
	if err != nil {
		panic(err)
	}
	return file
}

// File 多文件上传
func (r *Request) Files(name string) []*multipart.FileHeader {
	form, err := r.c.echo.MultipartForm()
	if err != nil {
		panic(err)
	}
	return form.File[name]
}

// RealIP 获取客户端真实IP
// 支持 X-Forwarded-For, RemoteAddr
func (r *Request) RealIP() string {
	return r.c.echo.RealIP()
}
func (r *Request) Url() string {
	return r.Http().URL.Path
}

// BindQuery 查询参数绑定
//
// tag: query
//
//	type T struct {
//		Id int `query:"id"`
//	    OrderRequest // 排序字段
//		PaginateRequest  // 分页字段
//	}
func (r *Request) BindQuery[T any]() T {
	var v T
	if err := echo.BindQueryParams(r.c.echo, &v); err != nil {
		panic(errors.NewBusinessf("参数解析错误 %v", err))
	}
	return v
}

// BindBody 请求体参数绑定
//
//	type T struct {
//		Id int `json:"id"`
//	    OrderRequest // 排序字段
//		PaginateRequest  // 分页字段
//	}
//
// tag: json/xml/form
func (r *Request) BindBody[T any]() T {
	var v T
	if err := echo.BindBody(r.c.echo, &v); err != nil {
		panic(errors.NewBusinessf("参数解析错误 %v", err))
	}
	return v
}

// Recorder 响应记录器
func (r *Request) Recorder() {
	if len(r.requestParam) > 0 {
		return
	}
	req := r.Http()

	contentType := req.Header.Get(echo.HeaderContentType)

	// uri
	r.requestParam["path"] = req.URL.Path
	r.requestParam["method"] = req.Method
	if strings.HasPrefix(contentType, echo.MIMEMultipartForm) {
		r.requestParam["contentType"] = "multipart/form-data"
	} else {
		r.requestParam["contentType"] = contentType
	}

	// query
	query := r.Querys()
	if len(query) > 0 {
		tempQuery := make(map[string]any, len(query))
		for k, v := range query {
			if len(v) == 1 {
				tempQuery[k] = v[0]
				return
			}
			tempQuery[k] = v
		}
		r.requestParam["query"] = tempQuery
	}
	// body
	body := getRequestBody(r.c)
	if body != "" {
		r.requestParam["body"] = body
	}
}

func (r *Request) GetRequestParam() map[string]any {
	return r.requestParam
}

func getRequestBody(c *Context) string {
	req := c.Request.Http()
	contentType := req.Header.Get(echo.HeaderContentType)
	headers := []string{
		echo.MIMEApplicationForm,
		echo.MIMEApplicationJSON,
		echo.MIMETextXML,
		echo.MIMEApplicationXML,
	}

	is := false
	for _, header := range headers {
		if strings.HasPrefix(contentType, header) {
			is = true
		}
	}
	if is == false {
		return ""
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	if bodyBytes == nil {
		return ""
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if strings.HasPrefix(contentType, echo.MIMEApplicationForm) {
		value := c.Request.Forms()
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return value.Encode()
	}
	return string(bodyBytes)
}

type PaginateRequest struct {
	Page     uint `query:"page" form:"page" json:"page"`
	PageSize uint `query:"pageSize" form:"pageSize" json:"pageSize"`
	OrderRequest
}

func (p PaginateRequest) GetPage() uint {
	return max(p.Page, 1)
}

func (p PaginateRequest) GetPageSize() uint {
	return max(cmp.Or(p.PageSize, 50), 1)
}

type OrderRequest struct {
	Order     string `query:"order" form:"order" json:"order"`
	Direction string `query:"direction" form:"direction" json:"direction"`
}

func (o OrderRequest) GetOrder() string {
	return o.Order
}

func (o OrderRequest) GetDirection() string {
	return o.Direction
}

// 获取请求参数
//
// ty
// query: Query参数(?id=)
// path: 路径参数(/:id)
// form: 表单参数
// "": Query参数(?id=) > 路径参数(/:id) > 表单参数
func param[T comparable](r *Request, ty, name string, defaults ...T) T {
	getters := []func(*echo.Context, string, ...any) (T, error){}
	switch ty {
	case "":
		getters = append(getters, echo.QueryParam[T], echo.PathParam[T], echo.FormValue[T])
	case "query":
		getters = append(getters, echo.QueryParam[T])
	case "path":
		getters = append(getters, echo.PathParam[T])
	case "form":
		getters = append(getters, echo.FormValue[T])
	default:
		panic(errors.NewBusinessf("参数类型 [%s] 不支持", ty))
	}

	var zero T
	var v T
	var err error
	for _, getter := range getters {
		v, err = getter(r.c.echo, name)
		if err == nil && v == zero {
			// 传递0值 相当于不传
			err = echo.ErrNonExistentKey
		}
		if errors.Is(err, echo.ErrNonExistentKey) {
			continue
		}
		if v != zero {
			return v
		}
	}
	if err != nil {
		if errors.Is(err, echo.ErrNonExistentKey) {
			if len(defaults) > 0 {
				return defaults[0]
			}
			panic(errors.NewBusinessf("参数 [%s] 必填", name))
		}
		panic(errors.NewBusinessf("参数 [%s] 输入错误: %v", name, err))
	}
	return v
}
