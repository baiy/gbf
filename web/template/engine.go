package template

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/web/template/funcs"
	"github.com/fsnotify/fsnotify"
)

// Engine 管理所有主题的模板缓存，支持热更新。
// 主题目录里面必须存在 template.md 文件
// 模板固定后缀
type Engine struct {
	mu           sync.Mutex
	driver       Driver
	DefaultTheme string            // 默认主题名称
	themeMap     map[string]*Theme // 缓存：主题名 -> 已编译模板集

	// 模板函数
	templateFunc funcs.FuncMap

	// 模板输出之前对内容最终的处理
	displayOutputBefore func([]byte) ([]byte, error)

	// 变量分隔符
	leftDelim  string
	rightDelim string

	// 模板变化监控
	watcher    *fsnotify.Watcher
	debounce   *time.Timer
	debounceMu sync.Mutex

	Logger log.Logger
}

// Themes 所有主题列表
func (e *Engine) Themes() []string {
	themes, _ := e.driver.Themes()
	// 默认主题移动到头部
	return append([]string{e.DefaultTheme}, slices.DeleteFunc(themes, func(s string) bool {
		return s == e.DefaultTheme
	})...)
}

// ThemeAddDefault 单一主题附加默认主题
func (e *Engine) ThemeAddDefault(theme string) []string {
	themes := []string{theme}
	if !slices.Contains(themes, e.DefaultTheme) {
		themes = append(themes, e.DefaultTheme)
	}
	return themes
}

// Fetch 模板渲染
func (e *Engine) Fetch(wr io.Writer, theme string, name string, data any) error {
	theme = cmp.Or(theme, e.DefaultTheme)
	t, ok := e.themeMap[theme]
	if !ok {
		return fmt.Errorf("unknown theme %s", theme)
	}
	w := &writer{
		wr:      wr,
		handler: e.displayOutputBefore,
	}
	err := t.fetch(w, name, data)
	if err != nil {
		return err
	}
	err = w.handle()
	if err != nil {
		return err
	}
	return nil
}

// Exists 模板是否存在
func (e *Engine) Exists(theme string, name string) bool {
	all, _ := e.driver.All(theme)
	return slices.Contains(all, name)
}

func (e *Engine) Delims(left, right string) *Engine {
	e.leftDelim = left
	e.rightDelim = right
	return e
}

// Build 主题编译
func (e *Engine) Build() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, theme := range e.Themes() {
		if t, ok := e.themeMap[theme]; ok {
			e.Logger.Debugf("build theme: %s", theme)
			err := t.build()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) DisplayOutputBefore(fc func([]byte) ([]byte, error)) *Engine {
	e.displayOutputBefore = fc
	return e
}

func NewEngine(driver Driver, defaultTheme string, logger log.Logger, templateFunc funcs.FuncMap) (*Engine, error) {
	engine := &Engine{
		DefaultTheme: cmp.Or(defaultTheme, "default"),
		driver:       driver,
		Logger:       logger,
		templateFunc: map[string]any{},
	}
	maps.Copy(engine.templateFunc, funcs.DefaultFuncMap)
	maps.Copy(engine.templateFunc, templateFunc)

	driver.Engine(engine)

	dth, err := driver.Themes()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(dth, engine.DefaultTheme) {
		return nil, errors.New("default theme not found")
	}

	themeMap := map[string]*Theme{}
	for _, theme := range engine.Themes() {
		themeMap[theme] = NewTheme(engine, theme)
	}
	engine.themeMap = themeMap

	err = engine.Build()
	if err != nil {
		return nil, err
	}

	err = engine.driver.Watch()
	if err != nil {
		return nil, err
	}
	return engine, nil
}
