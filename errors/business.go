package errors

import (
	"errors"
	"fmt"

	"github.com/baiy/gbf/util"
)

// Business 业务异常
type Business struct {
	Code    int
	Message string
	Data    any
	Url     string
}

func (e *Business) Error() string {
	return e.Message
}

// NewBusiness 业务异常
func NewBusiness(message string, options ...util.WithOption[Business]) *Business {
	b := &Business{Message: message}
	util.ApplyOption(b, options...)
	return b
}

// NewBusinessf 业务异常
//
//fmt:print
func NewBusinessf(message string, args ...any) *Business {
	return &Business{Message: fmt.Sprintf(message, args...)}
}

func IsBusiness(err error) *Business {
	v, ok := errors.AsType[*Business](err)
	if !ok {
		return nil
	}
	return v
}

func WithCode(code int) util.WithOption[Business] {
	return func(b *Business) {
		b.Code = code
	}
}

func WithData(data any) util.WithOption[Business] {
	return func(b *Business) {
		b.Data = data
	}
}

func WithUrl(url string) util.WithOption[Business] {
	return func(b *Business) {
		b.Url = url
	}
}
