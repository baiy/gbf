package gbf

import (
	"os"
	"strings"

	"github.com/baiy/gbf/util/yaml"
)

// LoadConfig 加载配置
//
// value string 文件路径
// value []byte 文件内容
func LoadConfig[T any](value any) T {
	var bc []byte
	if filePath, ok := value.(string); ok {
		c, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		bc = c
	}
	if fileContent, ok := value.([]byte); ok {
		bc = fileContent
	}

	if len(bc) < 1 {
		panic("empty config content")
	}

	var cfg T
	content := os.Expand(string(bc), func(varName string) string {
		// 支持环境变量默认值语法 ${VAR:-default}
		if strings.Contains(varName, ":-") {
			parts := strings.SplitN(varName, ":-", 2)
			envVal := os.Getenv(parts[0])
			if envVal == "" {
				return parts[1]
			}
			return envVal
		}
		return os.Getenv(varName)
	})
	err := yaml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
