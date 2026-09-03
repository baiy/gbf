package template

import (
	"html/template"
	"io"

	"github.com/baiy/gbf/log"
)

type Theme struct {
	name     string
	template *template.Template
	engine   *Engine
	logger   log.Logger
}

func NewTheme(engine *Engine, name string) *Theme {
	return &Theme{
		name:   name,
		engine: engine,
	}
}

func (t *Theme) fetch(wr io.Writer, name string, data any) error {
	return t.template.ExecuteTemplate(wr, name, data)
}

// 主题编译
func (t *Theme) build() error {
	bt := template.New("template_" + t.name)
	if t.engine.leftDelim != "" && t.engine.rightDelim != "" {
		bt.Delims(t.engine.leftDelim, t.engine.rightDelim)
	}
	names, err := t.engine.driver.All(t.name)
	t.engine.Logger.Debugf("build theme name: %v", names)
	if err != nil {
		return err
	}
	if len(t.engine.templateFunc) > 0 {
		bt.Funcs(t.engine.templateFunc)
	}
	for _, name := range names {
		b, err := t.engine.driver.Get(t.name, name)
		if err != nil {
			return err
		}
		t.engine.Logger.Infof("load template %s:%s", t.name, name)
		_, err = bt.New(name).Parse(string(b))
		if err != nil {
			return err
		}
	}
	t.template = bt
	return nil
}
