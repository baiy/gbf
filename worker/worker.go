package worker

import (
	"fmt"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util"
	"github.com/panjf2000/ants/v2"
)

type PanicHandler = func(str string)

// Worker 协程池
// 阻塞模式(默认): 提交任务会阻塞等待，直到池中有空闲的 Worker 来接手任务
// 非阻塞模式: 提交任务会立即返回，不会阻塞。如果池已满，则会直接返回一个错误
type Worker struct {
	name string
	ants *ants.Pool
	// 是否非阻塞模式 默认: false(阻塞)
	nonBlocking bool
	// 池大小 默认: 1000
	// CPU 密集型: CPU 核数 × 2（CPU 密集型）  例如: 计算/加密/压缩/图片处理
	// I/O 密集型: 1000-10000（I/O 密集型） 例如: 网络请求/数据库读写/文件操作
	size int

	logger Logger
}

func (w *Worker) Submit(task func()) {
	trackId := log.TrackId()
	err := w.ants.Submit(func() {
		// 捕获 panic
		err := errors.PanicToError(func() error {
			log.TrackReset(trackId)
			task()
			return nil
		})
		if err != nil {
			w.logger.Errorf("worker task [%s] error: %s", w.name, err)
		}
	})
	if err != nil {
		w.logger.Errorf("worker pool submit error: %s", err)
	}
}

func (w *Worker) Close() error {
	if w.ants != nil {
		w.ants.Release()
	}
	return nil
}

func New(options ...util.WithOption[Worker]) (*Worker, error) {
	w := &Worker{
		name:        time.Now().Format("20060102150405"),
		size:        1000,
		nonBlocking: false,
	}
	util.ApplyOption(w, options...)
	logger := Logger{logger: log.With(fmt.Sprintf("worker-%s", w.name))}
	w.logger = logger

	antsPool, err := ants.NewPool(
		max(w.size, 1),
		ants.WithLogger(logger),
		ants.WithExpiryDuration(time.Minute),
		ants.WithNonblocking(w.nonBlocking),
		ants.WithPanicHandler(func(str any) {
			logger.Errorf("worker panic: %s", str)
		}),
	)
	if err != nil {
		return nil, err
	}
	w.ants = antsPool
	logger.Infof("create worker size: %d blocking: %t", w.size, !w.nonBlocking)
	closer.Register(fmt.Sprintf("worker:%s", w.name), w)
	return w, nil
}

func WithName(name string) util.WithOption[Worker] {
	return func(p *Worker) {
		p.name = name
	}
}

func WithNonBlocking() util.WithOption[Worker] {
	return func(p *Worker) {
		p.nonBlocking = true
	}
}

func WithSize(size int) util.WithOption[Worker] {
	return func(p *Worker) {
		p.size = size
	}
}
