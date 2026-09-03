package closer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"
)

var container []*Item
var mu sync.Mutex
var CloseCtx, closeCtxCancel = context.WithCancel(context.Background())

type Item struct {
	Name  string
	Value Closer
}

type Closer interface {
	Close() error
}

func Register(name string, closer Closer) {
	mu.Lock()
	defer mu.Unlock()
	container = append(container, &Item{name, closer})
}

func GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info(fmt.Sprintf("system close start"))
	done := make(chan struct{}, 1)
	closeCtxCancel()
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		// 倒序关闭，后注册的资源先关闭
		for _, item := range slices.Backward(container) {
			slog.Info(fmt.Sprintf("system close [%s] start", item.Name))
			err := item.Value.Close()
			if err != nil {
				slog.Error(fmt.Sprintf("system close [%s] error:%s", item.Name, err))
			}
			slog.Info(fmt.Sprintf("system close [%s] end", item.Name))
		}
	}()
	select {
	case <-done:
		slog.Info("system close success")
		return
	case <-time.After(time.Second * 30):
		slog.Error("system close timeout, force exit")
		return
	}
}
