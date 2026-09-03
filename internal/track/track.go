package track

import (
	"strings"
	"uuid"

	"github.com/timandy/routine"
)

type Track struct {
	Id     string
	no     int64
	Source string
}

var threadLocal = routine.NewThreadLocal[*Track]()

func (t *Track) No() int64 {
	t.no = t.no + 1
	return t.no
}

func (t *Track) Reset(source string) {
	t.Id = id()
	t.no = 0
	t.Source = source
}

func Get() *Track {
	item := threadLocal.Get()
	if item == nil {
		item = &Track{
			Id: id(),
		}
		threadLocal.Set(item)
	}
	return item
}

func id() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
