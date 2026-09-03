package captcha

import (
	"github.com/baiy/gbf/log"
	"github.com/mojocn/base64Captcha"
)

const (
	numbers  = "0123456789"
	alphabet = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
)

type Captcha struct {
	captcha *base64Captcha.Captcha
	option  Option

	logger log.Logger
}

func (c *Captcha) Generate() (id, b64s, answer string, err error) {
	id, b64s, answer, err = c.captcha.Generate()
	if err != nil {
		c.logger.Errorf("create captcha error: %v", err)
		return
	}
	return
}

func (c *Captcha) Verify(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return c.captcha.Verify(id, answer, true)
}

func newCaptcha(key string, logger log.Logger, option *Option) *Captcha {
	var driver base64Captcha.Driver
	if option.style == "math" {
		driver = base64Captcha.NewDriverMath(
			option.height,
			option.width,
			2,
			0,
			nil,
			nil,
			nil,
		)
	} else {
		driver = base64Captcha.NewDriverString(
			option.height,
			option.width,
			2,
			0,
			option.length,
			numbers+alphabet,
			nil,
			nil,
			nil,
		)
	}

	option.store.SetTtl(option.ttl)
	option.store.SetLogger(logger)
	option.store.SetKey(key)
	return &Captcha{
		captcha: base64Captcha.NewCaptcha(driver, option.store),
		logger:  logger,
	}
}
