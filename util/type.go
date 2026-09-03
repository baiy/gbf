package util

type WithOption[T any] func(*T)

func ApplyOption[T any](value *T, options ...WithOption[T]) {
	for _, option := range options {
		option(value)
	}
}
