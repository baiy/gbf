package starter

import (
	"fmt"
	"slices"

	"github.com/baiy/gbf/closer"
)

type App struct {
	modules map[string]Starter // 模块初始化
	afters  []Action           // 初始化完成后执行队列

	sequence []string
}

// Sequence 添加初始化顺序
// name 添加的名称
// beforeName 添加在beforeName之后
func (a *App) Sequence(name string, beforeName string) *App {
	if slices.Contains(a.sequence, name) {
		panic(fmt.Sprintf("starter sequence exist: %s", name))
	}
	if beforeName != "" && !slices.Contains(a.sequence, beforeName) {
		panic(fmt.Sprintf("starter sequence not exist: %s", beforeName))
	}
	if beforeName == "" {
		a.sequence = append(a.sequence, name)
	} else {
		index := slices.Index(a.sequence, beforeName) + 1
		a.sequence = append(
			a.sequence[:index],
			append([]string{name}, a.sequence[index:]...)...)
	}

	return a
}

func (a *App) Add(name string, starter Starter, afters ...Action) *App {
	if !slices.Contains(a.sequence, name) {
		panic(fmt.Sprintf("starter sequence error: %s", name))
	}
	a.modules[name] = starter
	a.AppendAfter(afters...)
	return a
}

func (a *App) Timezone(starter Starter) *App {
	delete(a.modules, "timezone")
	return a.Add("timezone", starter)
}

func (a *App) Env(starter Starter) *App {
	delete(a.modules, "env")
	return a.Add("env", starter)
}

func (a *App) Log(starter Starter) *App {
	delete(a.modules, "log")
	return a.Add("log", starter)
}

func (a *App) Mysql(starter Starter, afters ...Action) *App {
	return a.Add("mysql", starter, afters...)
}

func (a *App) Rabbitmq(starter Starter, afters ...Action) *App {
	return a.Add("rabbitmq", starter, afters...)
}

func (a *App) Redis(starter Starter, afters ...Action) *App {
	return a.Add("redis", starter, afters...)
}

func (a *App) Cache(starter Starter, afters ...Action) *App {
	return a.Add("cache", starter, afters...)
}

func (a *App) Captcha(starter Starter, afters ...Action) *App {
	return a.Add("captcha", starter, afters...)
}

func (a *App) Lock(starter Starter, afters ...Action) *App {
	return a.Add("lock", starter, afters...)
}

func (a *App) Web(starter Starter, afters ...Action) *App {
	return a.Add("web", starter, afters...)
}

// AppendAfter 追加初始化完成后执行队列
func (a *App) AppendAfter(afters ...Action) *App {
	a.afters = append(a.afters, afters...)
	return a
}

// PrependAfter 插入初始化完成后执行队列
func (a *App) PrependAfter(afters ...Action) *App {
	a.afters = append(afters, a.afters...)
	return a
}

// Run 启动运行
// 会阻塞当前程序
func (a *App) Run() {
	// 初始化
	for _, name := range a.sequence {
		if module, ok := a.modules[name]; ok {
			after, err := module()
			if err != nil {
				panic(fmt.Sprintf("starter [%s] error: %v", name, err))
			}
			if after != nil {
				a.AppendAfter(after)
			}
		}
	}
	// 初始化完成执行
	for _, after := range a.afters {
		if err := after(); err != nil {
			panic(fmt.Sprintf("starter run error: %v", err))
		}
	}

	// 优雅关闭
	closer.GracefulShutdown()
}
