package captcha

import (
	"time"

	"github.com/baiy/gbf/log"
	"github.com/mojocn/base64Captcha"
)

type Store interface {
	base64Captcha.Store

	SetTtl(time.Duration)
	SetLogger(log.Logger)
	SetKey(string)
}

type Base struct {
	Key    string
	Ttl    time.Duration
	Logger log.Logger
}

func (b *Base) SetTtl(t time.Duration) {
	b.Ttl = t
}
func (b *Base) SetLogger(l log.Logger) {
	b.Logger = l

}
func (b *Base) SetKey(k string) {
	b.Key = k
}
