package web

import (
	"github.com/baiy/gbf/util/json"
	"github.com/labstack/echo/v5"
)

type JSONSerializer struct {
}

func (s *JSONSerializer) Serialize(c *echo.Context, target any, indent string) error {
	var v []byte
	var err error

	if indent == "" {
		v, err = json.Marshal(target)
	} else {
		v, err = json.Marshal(target, json.WithBeautify(indent))
	}
	if err != nil {
		return err
	}
	_, err = c.Response().Write(v)
	return err

	// 未使用 MarshalWrite 的原因: MarshalWrite 会多次调用写入方法, 造成After重复执行
	//if indent == "" {
	//	return json.MarshalWrite(c.Response(), target)
	//}
	//return json.MarshalWrite(c.Response(), target, json.WithBeautify(indent))
}

func (s *JSONSerializer) Deserialize(c *echo.Context, target any) error {
	if err := json.UnmarshalRead(c.Request().Body, target); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
