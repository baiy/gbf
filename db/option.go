package db

import (
	"github.com/baiy/gbf/util"
)

type Option struct {
	tablePrefix string // 表前缀
}

func WithTablePrefix(host string) util.WithOption[Option] {
	return func(o *Option) {
		o.tablePrefix = host
	}
}

func NewOption(options ...util.WithOption[Option]) *Option {
	option := &Option{}
	util.ApplyOption(option, options...)
	return option
}
