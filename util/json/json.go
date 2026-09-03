package json

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"

	"github.com/baiy/gbf/util"
)

type Option struct {
	indent       string
	indentPrefix string
}

func Marshal(in any, ops ...util.WithOption[Option]) ([]byte, error) {
	return json.Marshal(in, toJsonOption(ops)...)
}

func Unmarshal(in []byte, out any, ops ...util.WithOption[Option]) error {
	return json.Unmarshal(in, out, toJsonOption(ops)...)
}

func MarshalWrite(out io.Writer, in any, ops ...util.WithOption[Option]) error {
	return json.MarshalWrite(out, in, toJsonOption(ops)...)
}

func UnmarshalRead(in io.Reader, out any, ops ...util.WithOption[Option]) error {
	return json.UnmarshalRead(in, out, toJsonOption(ops)...)
}

func Encode(d any, ops ...util.WithOption[Option]) (string, error) {
	v, err := Marshal(d, ops...)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func Decode[T any](d string, ops ...util.WithOption[Option]) (T, error) {
	var v T
	err := Unmarshal([]byte(d), &v, ops...)
	if err != nil {
		return v, err
	}
	return v, err
}

func ToString(d any) string {
	v, _ := Marshal(d)
	return string(v)
}

type SerializeValue[T any] struct {
	V T
}

func Serialize[T any](value T, ops ...util.WithOption[Option]) string {
	data, err := Marshal(SerializeValue[T]{
		V: value,
	}, ops...)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Unserialize[T any](value string, ops ...util.WithOption[Option]) T {
	data := new(SerializeValue[T])
	err := Unmarshal([]byte(value), data, ops...)
	if err != nil {
		panic(err)
	}
	return data.V
}

func WithBeautify(ops ...string) util.WithOption[Option] {
	// ops[0] indent
	// ops[1] prefix
	if len(ops) == 0 {
		ops = []string{"  ", ""}
	}
	if len(ops) == 1 {
		ops = append(ops, "")
	}
	return func(o *Option) {
		o.indentPrefix = ops[0]
		o.indent = ops[1]
	}
}

func toJsonOption(options []util.WithOption[Option]) []json.Options {
	option := Option{}
	util.ApplyOption(&option, options...)

	var result []json.Options
	if option.indent != "" {
		result = append(result, jsontext.WithIndent(option.indent))
	}
	if option.indentPrefix != "" {
		result = append(result, jsontext.WithIndentPrefix(option.indentPrefix))
	}
	return result
}
