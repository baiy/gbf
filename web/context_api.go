package web

// api 请求
import (
	"cmp"
	"net/http"
	"strings"
	"time"

	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/util/yaml"
	"github.com/labstack/echo/v5"
)

// ApiHandlerFunc api 处理函数
//
// 返回值说明
// error 响应错误
// disableAutoResponse 禁用自动响应
// nil 响应成功 data:""
// value 响应成功+data:value
type ApiHandlerFunc func(c *ApiContext) any

const disableAutoResponse = "__disable_auto_response__"

type ApiContext struct {
	*Context
}

// DisableAutoResponse 禁用自动响应
// 处理逻辑中需要自己发送响应
func (c *ApiContext) DisableAutoResponse() any {
	return disableAutoResponse
}
func (c *ApiContext) SendResponse(r *ApiResponse) error {
	httpStatus := cmp.Or(r.HttpStatus, http.StatusOK)

	result := map[string]any{
		"code":    r.Code,
		"data":    r.Data,
		"message": r.Message,
		"status":  r.Status,
		"time":    time.Now().Format(time.RFC3339),
	}

	accept := c.Request.Header().Get("Accept")
	switch {
	// xml
	case strings.HasPrefix(accept, echo.MIMEApplicationXML):
		return c.Response.XML(result, httpStatus)
	// yaml
	case strings.HasPrefix(accept, "application/yaml"):
		b, err := yaml.Marshal(result)
		if err != nil {
			return err
		}
		return c.Response.Blob("application/yaml", b, httpStatus)
	// json
	default:
		return c.Response.JSON(result, httpStatus)
	}
}

func ApiHandlerWrapper(h ApiHandlerFunc) HandlerFunc {
	return func(c *Context) error {
		context := &ApiContext{Context: c}
		var r any
		err := errors.PanicToError(func() error {
			tempR := h(context)
			if err, ok := tempR.(error); ok {
				return err
			}
			r = tempR
			return nil
		})

		// 错误响应输出
		if err != nil {
			c.Logger.Errorf("api controller error:%v", err)
			if e := errors.IsBusiness(err); e != nil {
				return context.SendResponse(&ApiResponse{
					Code:    e.Code,
					Data:    e.Data,
					Status:  false,
					Message: e.Error(),
				})
			}
			data := ""
			if !env.IsProd() {
				data = err.Error()
			}
			return context.SendResponse(&ApiResponse{
				Code:    0,
				Data:    data,
				Message: "系统错误",
				Status:  false,
			})
		}

		// 不自动输出响应, 由控制器内部主动响应
		if rr, ok := r.(string); ok && rr == disableAutoResponse {
			return nil
		}

		// 返回值为nil时, 自动输出成功响应
		if r == nil {
			r = ""
		}

		// 自动输出成功响应
		return context.SendResponse(&ApiResponse{
			Data:    r,
			Status:  true,
			Message: "操作成功",
		})
	}
}

type ApiResponse struct {
	HttpStatus int
	Code       int
	Data       any
	Message    string
	Status     bool
}
