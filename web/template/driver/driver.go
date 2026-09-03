package driver

import (
	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/web/template"
)

// New 根据环境来确定使用那种模板文件驱动
// 开发环境使用本地文件系统 方便修改调试
// 其他环境 使用 embed.Fs 将模板嵌入二进制文件中
// f: NewFile
// e: NewEmbedFs
func New(f Create, e Create) Create {
	return func() template.Driver {
		if env.IsLocal() {
			return f()
		}
		return e()
	}
}
