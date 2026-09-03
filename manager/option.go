package manager

import "github.com/baiy/gbf/log"

type Option[O any] struct {
	Key    string
	Value  O
	Logger log.Logger
}
