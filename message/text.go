package message

import (
	"bytes"
	"strings"

	"github.com/baiy/gbf/log"
	"github.com/yuin/goldmark"
)

type Template struct {
	tpl    string
	values map[string]string
}

func (t *Template) AddValue(key string, v string) *Template {
	t.values[key] = v
	return t
}

func (t *Template) AddValues(values map[string]string) *Template {
	for k, v := range values {
		t.AddValue(k, v)
	}
	return t
}

func (t *Template) text() string {
	return t.String()
}

func (t *Template) MarkDownToHtml() string {
	md := t.String()
	source := []byte(md)
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		log.Errorf("markdown convert error:%s", err)
		return md
	}
	return buf.String()
}

func (t *Template) String() string {
	if len(t.values) == 0 {
		return t.tpl
	}
	v := t.tpl
	for key, value := range t.values {
		v = strings.ReplaceAll(v, "{"+key+"}", value)
	}
	return v
}

func NewTemplate(tpl string) *Template {
	return &Template{
		tpl:    tpl,
		values: map[string]string{},
	}
}
