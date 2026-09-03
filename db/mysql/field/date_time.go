package field

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/baiy/gbf/util/json"
	"github.com/baiy/gbf/util/times"
)

// DateTime mysql datetime/timestamp 数据类型
//
// sql
// `field`  datetime/timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "",
// `field`  datetime/timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "",
type DateTime string

func (d *DateTime) Scan(value any) (err error) {
	nullValue := &sql.NullTime{}
	err = nullValue.Scan(value)
	*d = NewDateTime(nullValue.Time)
	return
}

func (d DateTime) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d DateTime) GormDataType() string {
	return "datetime"
}

func (d DateTime) String() string {
	return string(d)
}

func (d DateTime) Time() time.Time {
	v, _ := time.Parse(time.DateTime, d.String())
	return v
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
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
	*d = NewDateTime(local)
	return nil
}

func NewDateTime(t time.Time) DateTime {
	return DateTime(t.Format(time.DateTime))
}
