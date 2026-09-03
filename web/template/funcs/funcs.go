package funcs

import (
	"html/template"
	"time"
)

type FuncMap = template.FuncMap

type ToTime interface {
	Time() time.Time
}

var DefaultFuncMap = FuncMap{
	"Time": Time, // 时间格式化
	"Cond": Cond, // 条件判断 类似三元表达式
}

func Time(t ToTime, formats ...string) string {
	if len(formats) == 0 {
		formats = []string{time.DateTime}
	}
	return t.Time().Format(formats[0])
}

func Cond(cond bool, trueVal, falseVal any) any {
	if cond {
		return trueVal
	}
	return falseVal
}
