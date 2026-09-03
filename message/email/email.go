package email

import (
	"fmt"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/worker"
	"github.com/wneessen/go-mail"
)

type Email struct {
	client *mail.Client
	logger log.Logger
}

func (e *Email) SendWithAsync(msg *Msg) {
	worker.Async(func() {
		err := e.send(msg)
		if err != nil {
			e.logger.Errorf("failed to send mail: %s", err)
		}
	})
}
func (e *Email) Send(msg *Msg) bool {
	err := e.send(msg)
	if err != nil {
		e.logger.Errorf("failed to send mail: %s", err)
		return false
	}
	return true
}

func (e *Email) send(msg *Msg) error {
	message := mail.NewMsg()
	if msg.from != "" {
		if err := message.From(msg.from); err != nil {
			return fmt.Errorf("failed to set From address: %s", err)
		}
	}
	if err := message.To(msg.to...); err != nil {
		return fmt.Errorf("failed to set To address: %s", err)
	}
	if err := message.Cc(msg.cc...); err != nil {
		return fmt.Errorf("failed to set Cc address: %s", err)
	}
	if err := message.Bcc(msg.bcc...); err != nil {
		return fmt.Errorf("failed to set Bcc address: %s", err)
	}
	for _, attach := range msg.attach {
		if err := message.AttachReader(attach.name, attach.file); err != nil {
			return fmt.Errorf("failed to attach file [%s]: %s", attach.name, err)
		}
	}
	message.Subject(msg.subject)
	message.SetBodyString(msg.contentType, msg.body)

	if err := e.client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}
	return nil
}
