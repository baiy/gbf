package client

import (
	"github.com/baiy/gbf/util/json"
)

type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type Response struct {
	body string

	err error
}

func (r *Response) Api[T any]() (ApiResponse[T], error) {
	if r.err != nil {
		return *new(ApiResponse[T]), r.err
	}
	return r.Json[ApiResponse[T]]()
}

func (r *Response) Json[T any]() (T, error) {
	if r.err != nil {
		return *new(T), r.err
	}

	return json.Decode[T](r.body)
}

func (r *Response) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.body, nil
}
