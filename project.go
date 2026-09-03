package gbf

import "github.com/baiy/gbf/util/yaml"

type Project struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// ParseProjectFile 解析项目信息 项目名称、版本号
//
// yaml文件内容 文件名一般为: gbf.yaml
// ```yaml
// # 项目名称
// name: app
// # 项目版本
// version: 0.0.0
// ```
func ParseProjectFile(yamlContent string) Project {
	var p Project
	err := yaml.Unmarshal([]byte(yamlContent), &p)
	if err != nil {
		panic(err)
	}
	return p
}
