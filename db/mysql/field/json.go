package field

import "gorm.io/datatypes"

type Json[T any] = datatypes.JSONType[T]
type JsonArray[T any] = datatypes.JSONSlice[T]

func NewJson[T any](data T) Json[T] {
	return datatypes.NewJSONType(data)
}

func NewJsonArray[T any](data []T) JsonArray[T] {
	return datatypes.NewJSONSlice(data)
}
