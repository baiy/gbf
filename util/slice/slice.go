package slice

import "github.com/duke-git/lancet/v2/slice"

func Map[T any, U any](sl []T, iteratee func(int, T) U) []U {
	return slice.Map(sl, iteratee)
}

func GroupWith[T any, U comparable](sl []T, iteratee func(item T) U) map[U][]T {
	return slice.GroupWith(sl, iteratee)
}

func Union[T comparable](slices ...[]T) []T {
	return slice.Union(slices...)
}

func Concat[T any](slices ...[]T) []T {
	return slice.Concat(slices...)
}

func FilterMap[T any, U any](sl []T, iteratee func(index int, item T) (U, bool)) []U {
	return slice.FilterMap(sl, iteratee)
}
