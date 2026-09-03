package template

import "io"

type writer struct {
	wr      io.Writer
	handler func([]byte) ([]byte, error)
	c       []byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	// 没有处理方法，直接写入
	if w.handler == nil {
		return w.wr.Write(p)
	}
	w.c = append(w.c, p...)
	return len(p), nil
}

func (w *writer) handle() error {
	if w.handler == nil {
		return nil
	}
	defer func() {
		w.c = nil
	}()

	// 执行处理方法
	r, err := w.handler(w.c)
	if err != nil {
		return err
	}
	_, err = w.wr.Write(r)
	if err != nil {
		return err
	}
	return nil
}
