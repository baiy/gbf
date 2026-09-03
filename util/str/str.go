package str

import (
	"strings"

	"github.com/baiy/gbf/util/json"
	"github.com/duke-git/lancet/v2/strutil"
)

func IsBlank(str string) bool {
	return strutil.IsBlank(str)
}

func Trim(str string, characterMask ...string) string {
	return strutil.Trim(str, characterMask...)
}

func IsNotBlank(str string) bool {
	return strutil.IsNotBlank(str)
}

func TemplateReplace(template string, data map[string]string) string {
	return strutil.TemplateReplace(template, data)
}

func Substring(str string, offset int, length uint) string {
	return strutil.Substring(str, offset, length)
}

func ToJsonString(d any) string {
	v, _ := json.Marshal(d)
	return string(v)
}

// Mask 通用脱敏函数
// s: 原始字符串
// keepLeft: 保留前几位
// keepRight: 保留后几位
// maskMaxLength: 掩码最大长度
func Mask(s string, keepLeft, keepRight int, maskMaxLength ...int) string {
	if s == "" {
		return ""
	}

	maskChar := "*"

	runes := []rune(s) // 支持中文等 Unicode
	length := len(runes)

	// 如果长度小于等于保留位数之和，直接返回原字符串（或全掩码，视业务而定）
	if length <= keepLeft+keepRight {
		return strings.Repeat(string(maskChar), length)
	}

	var builder strings.Builder
	builder.Grow(length)

	// 前半部分
	builder.WriteString(string(runes[:keepLeft]))
	// 中间掩码
	midLen := length - keepLeft - keepRight
	if len(maskMaxLength) > 0 {
		midLen = min(midLen, maskMaxLength[0])
	}
	builder.WriteString(strings.Repeat(maskChar, midLen))
	// 后半部分
	builder.WriteString(string(runes[length-keepRight:]))

	return builder.String()
}
