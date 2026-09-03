package captcha

import (
	"time"

	"github.com/baiy/gbf/manager"
)

func create(mo manager.Option[Option]) (*Captcha, error) {
	return newCaptcha(mo.Key, mo.Logger, &mo.Value), nil
}

var Manager = manager.New[*Captcha, Option](
	"captcha",
	create,
	nil,
	func() Option {
		return Option{
			ttl:    time.Minute * 10,
			height: 60,
			width:  240,
			style:  "math",
			length: 4,
		}
	},
)
