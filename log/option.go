package log

import (
	"log/slog"

	"github.com/baiy/gbf/util"
)

type Option struct {
	// 日志级别 DEBUG INFO WARN ERROR 默认:INFO
	level slog.Level

	// 日志是否输入到标准输出
	// 未启用文件日志的情况下 会强制启用标准输出
	stdout bool

	// 文件日志路径 为空不启用文件日志
	file string
	// 禁用文件日志滚动压缩
	fileRotationDisable bool
	// 文件日志保留数量 默认: 60
	fileMaxBackups int
	// 文件日志保留天数 默认: 365
	fileMaxAge int
}

func WithLevel(l slog.Level) util.WithOption[Option] {
	return func(o *Option) {
		o.level = l
	}
}

func WithDisableStdout() util.WithOption[Option] {
	return func(o *Option) {
		o.stdout = false
	}
}

func WithFile(path string) util.WithOption[Option] {
	return func(o *Option) {
		o.file = path
	}
}

func WithFileRotationDisable() util.WithOption[Option] {
	return func(o *Option) {
		o.fileRotationDisable = true
	}
}

func WithFileMaxBackups(n int) util.WithOption[Option] {
	return func(o *Option) {
		o.fileMaxBackups = n
	}
}

func WithFileMaxAge(n int) util.WithOption[Option] {
	return func(o *Option) {
		o.fileMaxAge = n
	}
}
