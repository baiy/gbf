// Package service 服务间调用
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/baiy/gbf/service/client"
)

type Service struct {
	// name 服务名称
	// 一般用于服务标识 k8s/docker中的服务名也必须这个值
	name string
	// port 服务端口
	// 容器内部监听监听端口
	port uint

	// 外部服务地址
	// 一般无需设置 name / port
	url string
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Port() uint {
	return s.port
}

func (s Service) Url() string {
	url := s.url
	if url == "" {
		url = fmt.Sprintf("http://%s:%d", s.name, s.port)
	}
	return strings.TrimSuffix(url, "/")
}

// Client 获取服务调用客户端
func (s Service) Client(ctxs ...context.Context) *client.Client {
	if len(ctxs) == 0 {
		client.New(s.Url(), context.Background())
	}
	return client.New(s.Url(), ctxs[0])
}

// NewService 创建内部服务
func NewService(name string, port uint) Service {
	return Service{
		name: name,
		port: port,
	}
}

// NewServiceByUrl 通过url创建外部服务
func NewServiceByUrl(url string) Service {
	return Service{
		url: url,
	}
}
