package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/baiy/gbf/internal/track"
)

const callerSkipSep = "_caller_skip_"

type Handler struct {
	slog.Handler
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{Handler: h.Handler.WithGroup(name)}
}
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	t := track.Get()
	// 追踪信息
	r.Add(slog.String("trace_id", t.Id))
	r.Add(slog.Int64("trace_no", t.No()))
	r.Add(slog.String("trace_source", t.Source))

	// 消息内容(调用栈处理)
	msg, call := handlerMessage(r.Message)
	r.Message = msg
	r.Add(slog.String("caller", call))

	return h.Handler.Handle(ctx, r)
}

func handlerMessage(msg string) (string, string) {
	// 为什么要自己计算调用栈, slog 和 自己封装的log 调用栈情况不太一样 这里磨平差异 统一输出
	var call string
	if strings.Index(msg, callerSkipSep) < 1 {
		// slog 原始调用
		call = caller(4)
	} else {
		// 自定义跳出层数
		temp := strings.Split(msg, callerSkipSep)
		msg = strings.Join(temp[1:], callerSkipSep)
		skip, err := strconv.Atoi(temp[0])
		// 自定义跳出层数异常时按照字符串处理
		if err != nil || skip < 0 || skip > 100 {
			call = temp[0]
		} else {
			call = caller(skip + 4)
		}
	}
	return msg, call
}

func caller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s %s:%d", file, funcName, line)
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		Handler: slog.NewJSONHandler(w, opts),
	}
}
