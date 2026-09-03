package field

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/baiy/gbf/util/json"
	"github.com/baiy/gbf/util/times"
)

// Date mysql date 数据类型
type Date string

func (d *Date) Scan(value any) (err error) {
	nullValue := &sql.NullTime{}
	err = nullValue.Scan(value)
	*d = NewDate(nullValue.Time)
	return
}

func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d Date) GormDataType() string {
	return "date"
}

func (d Date) String() string {
	return string(d)
}

func (d Date) Time() time.Time {
	v, _ := time.Parse(time.DateOnly, d.String())
	return v
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	if val == "" {
		return nil
	}

	local, err := times.Parse(val)
	if err != nil {
		return err
	}
	*d = NewDate(local)
	return nil
}

func NewDate(t time.Time) Date {
	return Date(t.Format(time.DateOnly))
}
