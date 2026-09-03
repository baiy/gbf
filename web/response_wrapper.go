package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type WrapperResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	size       int64
	statusCode int
}

func (w *WrapperResponseWriter) Write(b []byte) (n int, err error) {
	// 是否收集响应
	if w.isCollect() {
		w.body.Write(b)
	}
	n, err = w.ResponseWriter.Write(b)
	w.size += int64(n)
	return
}

func (w *WrapperResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w *WrapperResponseWriter) isCollect() bool {
	// 大于1M的响应不收集
	if w.size > 1024*1024 {
		return false
	}
	headers := []string{"text/html", "application/json"}

	contentType := w.Header().Get("Content-Type")

	for _, header := range headers {
		if strings.Contains(contentType, header) {
			return true
		}
	}
	return false
}

func (w *WrapperResponseWriter) String() string {
	body := w.body.String()
	if body == "" {
		return fmt.Sprintf(
			"response is empty or filtered, status:%v Content-Type:%v size:%v",
			w.statusCode,
			w.Header().Get("Content-Type"),
			w.size,
		)
	}
	return body
}

func (w *WrapperResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func NewWrapperResponseWriter(w http.ResponseWriter) *WrapperResponseWriter {
	return &WrapperResponseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}
