package times

import (
	"errors"
	"fmt"
	"time"
)

type parseValue interface {
	string | int64
}

var ParseTimeFormats = []string{
	time.DateTime, // 2006-01-02 15:04:05[.999999999]
	time.RFC3339,  // 2006-01-02T00:00:00[.999999999]Z / 2006-01-02T08:00:00[.999999999]+08:00 js: `(new Date()).toISOString()`
}

// Parse 时间解析 本地时区 time.Time
//
// 支持 时间字符串: ParseTimeFormats 时间戳: 范围如下
//
//	秒: 1900-01-01 00:00:00 - 2999-12-31 23:59:59
//	毫秒: (1900-01-01 00:00:00 - 1969-12-06 18:23:02) - (1971-01-12 12:47:31 - 2999-12-31 23:59:59)
//	微秒: (1900-01-01 00:00:00 - 1969-12-06 18:23:02) - (1971-01-12 12:47:31 - 2999-12-31 23:59:59)
func Parse[T parseValue](v T) (t time.Time, err error) {
	switch any(v).(type) {
	case string:
		return parseStr(any(v).(string))
	case int64:
		return parseUnix(any(v).(int64))
	default:
		return time.Time{}, fmt.Errorf("parse value error: %v", v)
	}
}
func parseStr(str string) (t time.Time, err error) {
	for _, format := range ParseTimeFormats {
		t, err = time.ParseInLocation(format, str, time.Local)
		if err == nil {
			if t.Location().String() != time.Local.String() {
				t = t.In(time.Local)
			}
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

func parseUnix(v int64) (time.Time, error) {
	var t time.Time

	if v >= -2209017943 && v <= 32503651199 {
		// 1900-01-01 00:00:00 - 2999-12-31 23:59:59
		t = time.Unix(v, 0)

	} else if v >= -2209017943000 && v <= 32503651199000 {
		// 毫秒
		// (1900-01-01 00:00:00 - 1969-12-06 18:23:02) - (1971-01-12 12:47:31 - 2999-12-31 23:59:59)
		t = time.UnixMilli(v)
	} else if v >= -2209017943000000 && v <= 32503651199000000 {
		// 微秒
		// (1900-01-01 00:00:00 - 1969-12-06 18:23:02) - (1971-01-12 12:47:31 - 2999-12-31 23:59:59)
		t = time.UnixMicro(v)
	}

	return t.In(time.Local), nil
}
