package mysql

import (
	"github.com/baiy/gbf/util"
)

type Option struct {
	// 地址
	Host string `yaml:"host"`
	// 端口
	Port string `yaml:"port"`
	// 用户名
	Username string `yaml:"username"`
	// 密码
	Password string `yaml:"password"`
	// 数据库名称
	Database string `yaml:"database"`

	// tls加密 https://github.com/go-sql-driver/mysql#tls
	Tls string `yaml:"tls"`

	// 连接池最大数量 默认值: 200
	MaxOpenConns uint `yaml:"maxOpenConns"`
	// 连接池最大空闲数量 默认值: 5
	MaxIdleConns uint `yaml:"maxIdleConns"`
	// 连接最大生命周期(秒) 默认值: 1小时
	ConnMaxLifetime uint `yaml:"connMaxLifetime"`
	// 连接最大空闲时间(秒) 默认值: 10分钟
	ConnMaxIdleTime uint `yaml:"connMaxIdleTime"`

	// 表前缀
	TablePrefix string `yaml:"tablePrefix"`
}

func WithHost(host string) util.WithOption[Option] {
	return func(o *Option) {
		o.Host = host
	}
}

func WithPort(port string) util.WithOption[Option] {
	return func(o *Option) {
		o.Port = port
	}
}

func WithUsername(username string) util.WithOption[Option] {
	return func(o *Option) {
		o.Username = username
	}
}

func WithPassword(password string) util.WithOption[Option] {
	return func(o *Option) {
		o.Password = password
	}
}

func WithDatabase(database string) util.WithOption[Option] {
	return func(o *Option) {
		o.Database = database
	}
}

func WithTls(tls string) util.WithOption[Option] {
	return func(o *Option) {
		o.Tls = tls
	}
}

func WithMaxOpenConns(maxOpenConns uint) util.WithOption[Option] {
	return func(o *Option) {
		o.MaxOpenConns = maxOpenConns
	}
}

func WithMaxIdleConns(maxIdleConns uint) util.WithOption[Option] {
	return func(o *Option) {
		o.MaxIdleConns = maxIdleConns
	}
}

func WithConnMaxLifetime(connMaxLifetime uint) util.WithOption[Option] {
	return func(o *Option) {
		o.ConnMaxLifetime = connMaxLifetime
	}
}

func WithConnMaxIdleTime(connMaxIdleTime uint) util.WithOption[Option] {
	return func(o *Option) {
		o.ConnMaxIdleTime = connMaxIdleTime
	}
}

func WithTablePrefix(tablePrefix string) util.WithOption[Option] {
	return func(o *Option) {
		o.TablePrefix = tablePrefix
	}
}
