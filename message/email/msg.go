package email

import (
	"io"
	"os"
	"path/filepath"

	"github.com/baiy/gbf/message"
	"github.com/wneessen/go-mail"
)

type Attach struct {
	name string
	file io.Reader
}

type Msg struct {
	from string   // 发件人 不设置为当前发送人邮箱
	to   []string // 收件人
	cc   []string // 抄送

	bcc []string // 密送

	attach []Attach // 附件

	subject     string
	body        string
	contentType mail.ContentType
}

func (m *Msg) To(email ...string) {
	m.to = append(m.to, email...)
}

func (m *Msg) Bcc(email ...string) {
	m.bcc = append(m.bcc, email...)
}

func (m *Msg) Cc(email ...string) {
	m.cc = append(m.cc, email...)
}

func (m *Msg) From(email string) {
	m.from = email
}

func (m *Msg) Subject(value *message.Template) {
	m.subject = value.String()
}

func (m *Msg) BodyWithTextPlain(value *message.Template) {
	m.contentType = mail.TypeTextPlain
	m.body = value.String()
}

func (m *Msg) BodyWithHtml(value *message.Template) {
	m.contentType = mail.TypeTextHTML
	m.body = value.String()
}

func (m *Msg) BodyWithHtmlByMarkdown(value *message.Template) {
	m.contentType = mail.TypeTextHTML
	m.body = value.MarkDownToHtml()
}

func (m *Msg) Attach(name string, file io.Reader) {
	m.attach = append(m.attach, Attach{
		name: name,
		file: file,
	})
}

func (m *Msg) AttachWithPath(path string) error {
	r, err := getFileReader(path)
	if err != nil {
		return err
	}
	m.Attach(filepath.Base(path), r)
	return nil
}

func getFileReader(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
