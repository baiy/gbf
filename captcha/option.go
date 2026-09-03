package captcha

import (
	"time"

	"github.com/baiy/gbf/util"
)

type Option struct {
	store  Store
	ttl    time.Duration // 过期时间
	width  int           // 宽度
	height int           // 高度
	length int           // 长度
	style  string        // 样式: string(数字字母), math(数学计算)
}

func WithStore(s Store) util.WithOption[Option] {
	return func(o *Option) {
		o.store = s
	}
}

func WithTtl(ttl time.Duration) util.WithOption[Option] {
	return func(o *Option) {
		o.ttl = ttl
	}
}

func WithHeight(height int) util.WithOption[Option] {
	return func(o *Option) {
		o.height = height
	}
}

func WithWidth(width int) util.WithOption[Option] {
	return func(o *Option) {
		o.width = width
	}
}

func WithLength(length int) util.WithOption[Option] {
	return func(o *Option) {
		o.length = length
	}
}

func WithMath() util.WithOption[Option] {
	return func(o *Option) {
		o.style = "math"
	}
}
