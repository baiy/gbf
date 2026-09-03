package field

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/baiy/gbf/util/json"
	"github.com/baiy/gbf/util/times"
)

// DateTimeMillisecond 毫秒精度
//
// sql
// `field`  datetime(3)/timestamp(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT "",
// `field`  datetime(3)/timestamp(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT "",
type DateTimeMillisecond string

func (d *DateTimeMillisecond) Scan(value any) (err error) {
	nullValue := &sql.NullTime{}
	err = nullValue.Scan(value)
	*d = NewDateTimeMillisecond(nullValue.Time)
	return
}

func (d DateTimeMillisecond) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d DateTimeMillisecond) GormDataType() string {
	return "datetime"
}

func (d DateTimeMillisecond) String() string {
	return string(d)
}

func (d DateTimeMillisecond) Time() time.Time {
	v, _ := time.Parse("2006-01-02 15:04:05.000", d.String())
	return v
}

func (d DateTimeMillisecond) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *DateTimeMillisecond) UnmarshalJSON(b []byte) error {
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
	*d = NewDateTimeMillisecond(local)
	return nil
}

func NewDateTimeMillisecond(t ...time.Time) DateTimeMillisecond {
	if len(t) == 0 {
		t = append(t, time.Now())
	}
	return DateTimeMillisecond(t[0].Format("2006-01-02 15:04:05.000"))
}
