package manager

type Instance[T any, O any] struct {
	Key    string
	Value  T
	Option *O
}

func newInstance[T any, O any](key string, value T, option *O) *Instance[T, O] {
	return &Instance[T, O]{
		Key:    key,
		Value:  value,
		Option: option,
	}
}
