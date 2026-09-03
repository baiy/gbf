package driver

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/web/template"
)

// EmbedFs embed.FS模板获取器 一般用于生产环境
type EmbedFs struct {
	embedFs embed.FS // 模板根目录 fs/rootDir 实际的模板目录
	rootDir string   // 这个路径是相对引入embed.FS的文件路径
	suffix  string   // 模板文件后缀 默认值: .tpl
	engine  *template.Engine

	logger log.Logger
}

func (e *EmbedFs) fs() fs.FS {
	sf, err := fs.Sub(e.embedFs, e.rootDir)
	if err != nil {
		panic(err)
	}
	return sf
}

func (e *EmbedFs) Get(theme string, name string) ([]byte, error) {
	for _, t := range e.engine.ThemeAddDefault(theme) {
		tplPath := filepath.Join(t, name) + e.suffix
		if _, err := fs.Stat(e.fs(), filepath.Join(t, name)+e.suffix); err == nil {
			b, err := fs.ReadFile(e.fs(), tplPath)
			if err == nil {
				return b, nil
			}
			e.logger.Errorf("reading file [%s] error: %s", tplPath, err)
		}
	}
	return nil, errors.New(fmt.Sprintf("template %s:%s not found", theme, name))
}

func (e *EmbedFs) All(theme string) ([]string, error) {
	var names []string
	for _, t := range e.engine.ThemeAddDefault(theme) {

		err := fs.WalkDir(e.fs(), t, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, e.suffix) {
				return nil
			}
			name := strings.TrimPrefix(
				strings.TrimSuffix(strings.TrimPrefix(path, t), e.suffix), string(filepath.Separator),
			)
			if name == "" {
				return errors.New("empty template name")
			}
			if !slices.Contains(names, name) {
				names = append(names, name)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return names, nil
}

func (e *EmbedFs) Themes() ([]string, error) {
	entries, err := fs.ReadDir(e.fs(), ".")
	if err != nil {
		e.logger.Errorf("reading dir [%s] error:%s", e.rootDir, err)
		return nil, err
	}
	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		theme := entry.Name()
		// 检查该子目录下是否存在 template.md 文件
		if _, err := fs.Stat(e.fs(), filepath.Join(theme, "template.md")); err != nil {
			e.logger.Warnf("theme %s is invalid: missing template.md", theme)
			continue
		}
		themes = append(themes, theme)
	}
	return themes, nil
}

func (e *EmbedFs) Engine(engine *template.Engine) {
	e.engine = engine
	e.logger = e.engine.Logger.With("embed")
}

func (e *EmbedFs) Watch() error {
	return nil
}

func NewEmbedFs(fs embed.FS, rootDir string, suffixs ...string) Create {
	suffix := ".tpl"
	if len(suffixs) > 0 && suffixs[0] != "" {
		suffix = suffixs[0]
	}

	return func() template.Driver {
		return &EmbedFs{embedFs: fs, rootDir: rootDir, suffix: suffix}
	}
}
