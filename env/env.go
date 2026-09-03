package env

import (
	"strings"
)

var env = Prod

const (
	Prod    = "PROD"    // 生产
	Local   = "LOCAL"   // 本地
	Dev     = "DEV"     // 开发
	Test    = "TEST"    // 测试
	Staging = "STAGING" // 预发布
	Uat     = "UAT"     // 客户验证
)

func parse(v string) string {
	v = strings.ToUpper(strings.TrimSpace(v))
	switch v {
	case "DEV":
		return Dev
	case "LOCAL":
		return Local
	case "TEST":
		return Test
	case "STAGING":
		return Staging
	case "UAT":
		return Uat
	case "PROD":
		return Prod
	}
	return Prod
}

func IsProd() bool {
	return env == Prod
}

func IsTest() bool {
	return env == Test
}

func IsUat() bool {
	return env == Uat
}

func IsStaging() bool {
	return env == Staging
}

func IsDev() bool {
	return env == Dev
}

func IsLocal() bool {
	return env == Local
}

// Set 设置环境
// v PROD / LOCAL / DEV / TEST / STAGING / UAT
func Set(v string) {
	env = parse(v)
}

func Get() string {
	return env
}
