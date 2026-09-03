package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/baiy/gbf/env"
)

// PanicToError 将 panic 转换为 error
func PanicToError(fn func() error) (err error) {
	defer func() {
		// 捕获 panic
		if r := recover(); r != nil {
			if !env.IsProd() {
				debug.PrintStack()
			}
			slog.Warn(fmt.Sprintf("panic: %s", debug.Stack()))
			tmpErr, ok := r.(error)
			if ok {
				err = tmpErr
				return
			}
			err = fmt.Errorf("%v", r)
		}
	}()
	return fn()
}

var Is = errors.Is
var As = errors.As
var Join = errors.Join

func AsType[E error](err error) (E, bool) {
	return errors.AsType[E](err)
}
