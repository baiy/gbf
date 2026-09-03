package gbf

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/starter"
	"github.com/baiy/gbf/util"
	"github.com/baiy/gbf/util/cryptor"
)

var app *starter.App

type AppOption struct {
	name    string
	version string
	tz      *time.Location
	env     func() string
	logFile string

	secret string
}

// NewApp 创建应用实例
func NewApp(projectYamlContent string, options ...util.WithOption[AppOption]) *starter.App {
	project := ParseProjectFile(projectYamlContent)
	tz, err := time.LoadLocation(cmp.Or(os.Getenv("TZ"), "Asia/Shanghai"))
	if err != nil {
		panic(err)
	}
	option := &AppOption{
		name:    project.Name,
		version: project.Version,
		tz:      tz,
		env: func() string {
			return cmp.Or(os.Getenv("RUNTIME_ENV"), env.Prod)
		},
		logFile: "/data/log/app.log",
	}
	util.ApplyOption(option, options...)
	constant.AppName = option.name
	constant.AppVersion = option.version
	constant.Secret = cmp.Or(option.secret, os.Getenv("APP_SECRET"))
	if constant.Secret == "" {
		slog.Warn(fmt.Sprintf("secret 没有自定义设置, 使用项目名称 %s 的 MD5 值作为 secret", option.name))
		constant.Secret = cryptor.Md5(option.name)
	}

	app = starter.NewApp().
		Timezone(tzStarter(option.tz)).
		Env(env.Starter(option.env())).
		Log(starter.Wrapper(func() error {
			var logOptions []util.WithOption[log.Option]
			if env.IsLocal() {
				// 本地开发环境强制日志级别为DEBUG 无文件输出
				logOptions = append(logOptions, log.WithFile(""), log.WithLevel(slog.LevelDebug))
			} else {
				logOptions = append(logOptions, log.WithFile(option.logFile))
			}
			return log.Init(logOptions...)
		}))

	return app
}

func App() *starter.App {
	return app
}

func WithTz(tz *time.Location) util.WithOption[AppOption] {
	return func(option *AppOption) {
		option.tz = tz
	}
}

func WithEnv(env func() string) util.WithOption[AppOption] {
	return func(option *AppOption) {
		option.env = env
	}
}

func WithLogFile(logFile string) util.WithOption[AppOption] {
	return func(option *AppOption) {
		option.logFile = logFile
	}
}

func WithSecret(secret string) util.WithOption[AppOption] {
	return func(option *AppOption) {
		option.secret = secret
	}
}

type InfoData struct {
	GbfVersion string
	AppVersion string
	AppName    string
}

func Info() InfoData {
	return InfoData{
		GbfVersion: constant.BaseVersion,
		AppVersion: constant.AppVersion,
		AppName:    constant.AppName,
	}
}

func tzStarter(tz *time.Location) starter.Starter {
	return starter.Wrapper(func() error {
		time.Local = tz
		return nil
	})
}
