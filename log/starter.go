package log

import (
	"github.com/baiy/gbf/starter"
	"github.com/baiy/gbf/util"
)

// Starter 创建全局启动器
func Starter(options ...util.WithOption[Option]) starter.Starter {
	return starter.Wrapper(func() error {
		return Init(options...)
	})
}
