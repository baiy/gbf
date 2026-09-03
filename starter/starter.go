package starter

type Action func() error

type Starter func() (after Action, err error)

// Wrapper Starter包装方法
// Wrapper(init)
// Wrapper(init,after)
func Wrapper(actions ...Action) Starter {
	return func() (Action, error) {
		action := actions[0]
		err := action()
		if err != nil {
			return nil, err
		}
		if len(actions) > 1 {
			return actions[1], nil
		}
		return nil, nil
	}
}

// 默认初始化顺序
var defaultSequence = []string{
	"timezone", "env", "log",
	"redis", "lock", "mysql", "rabbitmq",
	"cache", "captcha", "email",
	"web",
}

func NewApp() *App {
	return &App{
		modules:  make(map[string]Starter),
		sequence: defaultSequence,
	}
}
