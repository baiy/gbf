package log

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/DeRuina/timberjack"
	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/env"
	"github.com/baiy/gbf/internal/track"
	"github.com/baiy/gbf/util"
	"github.com/baiy/gbf/util/str"
)

var baseLogger *slog.Logger
var logger Logger

const callerSkipSep = "_caller_skip_"

type Logger struct {
	channel string
	slogger *slog.Logger
}

func (l Logger) Slog() *slog.Logger {
	return l.slogger
}

func (l Logger) Log(skip int, level slog.Level, msg string, args ...any) {
	// 附加调用来源栈跳出层数 方便后续计算调用来源
	msg = fmt.Sprintf("%d%s%s", skip+1, callerSkipSep, msg)
	l.slogger.Log(context.Background(), level, msg, args...)
}

func (l Logger) Logf(skip int, level slog.Level, msg string, args ...any) {
	l.Log(skip+1, level, fmt.Sprintf(msg, args...))
}

func (l Logger) Print(level slog.Level, msg string, args ...any) {
	l.Log(1, level, msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.Log(1, slog.LevelError, msg, args...)
}
func (l Logger) Info(msg string, args ...any) {
	l.Log(1, slog.LevelInfo, msg, args...)
}
func (l Logger) Debug(msg string, args ...any) {
	l.Log(1, slog.LevelDebug, msg, args...)
}
func (l Logger) Warn(msg string, args ...any) {
	l.Log(1, slog.LevelWarn, msg, args...)
}

func (l Logger) Printf(level slog.Level, msg string, args ...any) {
	l.Log(1, level, fmt.Sprintf(msg, args...))
}

//fmt:print
func (l Logger) Errorf(msg string, args ...any) {
	l.Log(1, slog.LevelError, fmt.Sprintf(msg, args...))
}

//fmt:print
func (l Logger) Infof(msg string, args ...any) {
	l.Log(1, slog.LevelInfo, fmt.Sprintf(msg, args...))
}

//fmt:print
func (l Logger) Debugf(msg string, args ...any) {
	l.Log(1, slog.LevelDebug, fmt.Sprintf(msg, args...))
}

//fmt:print
func (l Logger) Warnf(msg string, args ...any) {
	l.Log(1, slog.LevelWarn, fmt.Sprintf(msg, args...))
}

func (l Logger) With(channel string) Logger {
	if channel == "" {
		return l
	}
	if l.channel == "global" {
		return newLogger(channel)
	}
	return newLogger(fmt.Sprintf("%s.%s", l.channel, channel))
}

// Get 获取日志默认实例
func Get() Logger {
	return logger
}

func With(channel string) Logger {
	return logger.With(channel)
}

func Error(msg string, args ...any) {
	logger.Log(1, slog.LevelError, msg, args...)
}
func Info(msg string, args ...any) {
	logger.Log(1, slog.LevelInfo, msg, args...)
}
func Debug(msg string, args ...any) {
	logger.Log(1, slog.LevelDebug, msg, args...)
}
func Warn(msg string, args ...any) {
	logger.Log(1, slog.LevelWarn, msg, args...)
}

//fmt:print
func Errorf(msg string, args ...any) {
	logger.Log(1, slog.LevelError, fmt.Sprintf(msg, args...))
}

//fmt:print
func Infof(msg string, args ...any) {
	logger.Log(1, slog.LevelInfo, fmt.Sprintf(msg, args...))
}

//fmt:print
func Debugf(msg string, args ...any) {
	logger.Log(1, slog.LevelDebug, fmt.Sprintf(msg, args...))
}

//fmt:print
func Warnf(msg string, args ...any) {
	logger.Log(1, slog.LevelWarn, fmt.Sprintf(msg, args...))
}

func TrackReset(sources ...string) {
	source := ""
	if len(sources) > 0 {
		source = sources[0]
	}
	track.Get().Reset(source)
}

func TrackId() string {
	return track.Get().Id
}

// ElapsedTime 运行耗时
// 在方法首行使用 defer ElapsedTime("标识符")()
func ElapsedTime(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		slog.Info(fmt.Sprintf("[%s] elapsed time: %.3f(s) ", name, elapsed.Seconds()))
	}
}

func Init(options ...util.WithOption[Option]) error {
	option := &Option{
		level:          slog.LevelInfo,
		stdout:         true,
		fileMaxBackups: 60,
		fileMaxAge:     365,
	}
	util.ApplyOption(option, options...)

	baseLogger = slog.New(slog.NewJSONHandler(writer(option), &slog.HandlerOptions{
		Level:       option.level,
		AddSource:   false,
		ReplaceAttr: replaceAttr,
	}))
	slog.SetDefault(baseLogger.With("channel", "global").WithGroup("attr"))

	// 设置默认日志实例
	logger = newLogger("global")
	Infof("[%s:%s gbf:%s ]:%s log init secret:%s option: %+v",
		constant.AppName,
		constant.AppVersion,
		constant.BaseVersion,
		env.Get(),
		str.Mask(constant.Secret, 12, 12),
		option,
	)
	return nil
}

func newLogger(channel string) Logger {
	return Logger{
		channel: channel,
		slogger: baseLogger.With("channel", channel).WithGroup("attr"),
	}
}

func writer(o *Option) io.Writer {
	var writers []io.Writer

	if o.file != "" {
		if o.fileRotationDisable {
			file, err := os.OpenFile(o.file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(fmt.Sprintf("打开文件[%s]失败: %s", o.file, err))
			}
			writers = append(writers, file)
		} else {
			// 日志滚动压缩
			writers = append(writers, &timberjack.Logger{
				Filename:         o.file,
				MaxBackups:       cmp.Or(o.fileMaxBackups, 60),
				MaxSize:          0,
				MaxAge:           cmp.Or(o.fileMaxAge, 365),
				Compression:      "gzip",
				RotateAt:         []string{"00:00"},
				LocalTime:        true,
				BackupTimeFormat: "2006-01-02-15-04-05",
			})
		}
	}
	if o.stdout || len(writers) == 0 {
		writers = append([]io.Writer{os.Stdout}, writers...)
	}
	return io.MultiWriter(writers...)
}

// ReplaceAttr 替换属性
// 实现track/调用栈功能
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Any(slog.TimeKey, time.Now().Format("2006-01-02T15:04:05.999999Z07:00"))
	}
	if a.Key == slog.MessageKey {
		return slog.Any(slog.MessageKey, buildMessage(a.Value.String()))
	}
	return a
}

func buildMessage(msg string) map[string]any {
	message := map[string]any{}

	// track 信息
	t := track.Get()
	trackValue := map[string]any{}
	trackValue["id"] = t.Id
	trackValue["no"] = t.No()
	if t.Source != "" {
		trackValue["source"] = t.Source
	}
	message["track"] = trackValue

	// 为什么要自己计算调用栈, slog 和 自己封装的log 调用栈情况不太一样 这里磨平差异 统一输出
	var call slog.Source
	if strings.Index(msg, callerSkipSep) < 1 {
		// slog 原始调用
		call = caller(7)
	} else {
		// 自定义跳出层数
		temp := strings.Split(msg, callerSkipSep)
		skip, err := strconv.Atoi(temp[0])
		if err != nil {
			call = caller(7)
			msg = fmt.Sprintf("log write error: %s | %s", err, msg)
		} else {
			call = caller(skip + 7)
			msg = strings.Join(temp[1:], callerSkipSep)
		}
	}
	message["caller"] = call
	message["content"] = msg
	return message
}

func caller(skip int) slog.Source {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return slog.Source{}
	}
	funcName := runtime.FuncForPC(pc).Name()
	return slog.Source{
		Function: funcName,
		File:     file,
		Line:     line,
	}
}
