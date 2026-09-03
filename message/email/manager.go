package email

import (
	"github.com/baiy/gbf/manager"
	"github.com/wneessen/go-mail"
)

func create(mo manager.Option[Option]) (*Email, error) {
	option := mo.Value
	client, err := mail.NewClient(
		option.host,
		append(
			[]mail.Option{
				mail.WithUsername(option.username),
				mail.WithPassword(option.password),
				mail.WithSMTPAuth(option.smtpAuthType),
				mail.WithLogger(Logger{logger: mo.Logger}),
			},
			option.mailOption...,
		)...,
	)
	if err != nil {
		return nil, err
	}
	return &Email{
		client: client,
		logger: mo.Logger,
	}, nil
}

var Manager = manager.New[*Email, Option](
	"email",
	create,
	nil,
	func() Option {
		return Option{smtpAuthType: mail.SMTPAuthLogin}
	},
)
