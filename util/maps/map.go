package maps

import (
	"github.com/duke-git/lancet/v2/maputil"
)

func Keys[K comparable, V any](m map[K]V) []K {
	return maputil.Keys(m)
}

func Values[K comparable, V any](m map[K]V) []V {
	return maputil.Values(m)
}

func MapValues[K comparable, V any, T any](m map[K]V, iteratee func(key K, value V) T) map[K]T {
	return maputil.MapValues(m, iteratee)
}

func FilterByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V {
	return maputil.FilterByValues(m, values)
}

func FilterByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	return maputil.FilterByKeys(m, keys)
}
