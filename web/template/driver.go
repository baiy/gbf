package template

type Driver interface {
	// Get 获取模板内容
	Get(theme string, name string) ([]byte, error)
	// All 获取所有模板名称
	All(theme string) ([]string, error)
	// Themes 获取所有主题名称
	Themes() ([]string, error)
	// Engine 注入模板引擎对象
	Engine(*Engine)
	// Watch 监控模板变化
	Watch() error
}
