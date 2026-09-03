package driver

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/web/template"
	"github.com/fsnotify/fsnotify"
)

type Create func() template.Driver

// File 本地文件模板获取器 一般用于开发环境
type File struct {
	rootDir string // 模板根目录 相对于 os.Getwd()
	suffix  string // 模板文件后缀 默认值: .tpl
	engine  *template.Engine

	// 模板变化监控
	watcher    *fsnotify.Watcher
	debounce   *time.Timer
	debounceMu sync.Mutex

	logger log.Logger
}

func (f *File) Get(theme string, name string) ([]byte, error) {
	for _, t := range f.engine.ThemeAddDefault(theme) {
		tplPath := filepath.Join(f.rootDir, t, name) + f.suffix
		if _, err := os.Stat(tplPath); err == nil {
			b, err := os.ReadFile(tplPath)
			if err == nil {
				return b, nil
			}
			f.logger.Errorf("reading file [%s] error:%s", tplPath, err)
		}
	}
	return nil, errors.New(fmt.Sprintf("file template %s:%s not found", theme, name))
}

func (f *File) All(theme string) ([]string, error) {
	var names []string
	for _, t := range f.engine.ThemeAddDefault(theme) {
		themeDir := filepath.Join(f.rootDir, t)
		err := filepath.WalkDir(themeDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, f.suffix) {
				return nil
			}
			name := strings.TrimPrefix(
				strings.TrimSuffix(strings.TrimPrefix(path, themeDir), f.suffix), string(filepath.Separator),
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

func (f *File) Themes() ([]string, error) {
	entries, err := os.ReadDir(f.rootDir)
	if err != nil {
		return nil, err
	}
	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		theme := entry.Name()
		// 检查该子目录下是否存在 template.md 文件
		if _, err := os.Stat(filepath.Join(f.rootDir, theme, "template.md")); err != nil {
			f.logger.Warnf("theme %s is invalid: missing template.md", theme)
			continue
		}
		themes = append(themes, theme)
	}
	return themes, nil
}

func (f *File) Engine(engine *template.Engine) {
	f.engine = engine
	f.logger = f.engine.Logger.With("file")
}

func (f *File) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	f.watcher = watcher

	// 递归监控整个模板根目录
	err = filepath.WalkDir(f.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return f.fileWatcherAdd(path)
		}
		return nil
	})
	if err != nil {
		err1 := watcher.Close()
		if err1 != nil {
			return err1
		}
		return err
	}
	go f.watchLoop()
	closer.Register(fmt.Sprintf("web.controller.template.file:%s", f.rootDir), f)
	return nil
}

func (f *File) Close() error {
	return f.watcher.Close()
}

// fileWatch 模板文件变化监听
func (f *File) fileWatcherAdd(p string) error {
	f.logger.Infof("file watcher add: %s", p)
	return f.watcher.Add(p)
}

// watchLoop 文件监控事件循环，带防抖机制。
func (f *File) watchLoop() {
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}

			// 新建目录
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if err := f.fileWatcherAdd(event.Name); err != nil {
						f.logger.Warnf("file watcher add error:%s", err)
					}
					continue
				}
			}

			// 防抖
			f.debounceMu.Lock()
			if f.debounce != nil {
				f.debounce.Stop()
			}

			duration := time.Second * 10
			if env.IsLocal() {
				duration = time.Second * 2
			}
			f.debounce = time.AfterFunc(duration, func() {
				// 重新编译模板
				f.logger.Debugf("file watcher trigger build:%s", event.String())
				err := f.engine.Build()
				if err != nil {
					f.logger.Warnf("build error:%s", err)
				}
			})
			f.debounceMu.Unlock()

		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			f.logger.Errorf("watcher error:%s", err)
		}
	}
}

func NewFile(rootDir string, suffixs ...string) Create {
	suffix := ".tpl"
	if len(suffixs) > 0 && suffixs[0] != "" {
		suffix = suffixs[0]
	}
	return func() template.Driver {
		return &File{rootDir: rootDir, suffix: suffix}
	}
}
