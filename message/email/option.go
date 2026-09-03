package email

import "github.com/wneessen/go-mail"

type Option struct {
	host         string
	username     string
	password     string
	smtpAuthType mail.SMTPAuthType
	mailOption   []mail.Option
}
