package worker

import "github.com/baiy/gbf/errors"
import "github.com/baiy/gbf/log"

// Async 异步执行
//
// 封装 跟踪id 重置
func Async(fc func()) {
	trackId := log.TrackId()
	go func() {
		log.TrackReset(trackId)
		// 捕获 panic 防止程序崩溃
		err := errors.PanicToError(func() error {
			fc()
			return nil
		})
		if err != nil {
			log.Errorf("Async error: %v", err)
		}
	}()
}
