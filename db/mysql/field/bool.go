package field

import (
	"database/sql"
	"database/sql/driver"

	"github.com/baiy/gbf/util/json"
)

// Bool mysql tinyint 数据类型 1 是 0 否
//
// `field` tinyint unsigned NOT NULL DEFAULT '0' COMMENT ”,
type Bool bool

func (b Bool) Is() bool {
	return bool(b)
}

func (b Bool) Value() (driver.Value, error) {
	if b.Is() {
		return int64(1), nil
	}
	return int64(0), nil
}

func (b *Bool) Scan(value any) (err error) {
	nullInt64 := &sql.NullInt64{}
	err = nullInt64.Scan(value)
	*b = NewBool(nullInt64.Int64 == 1)
	return
}

func (b Bool) GormDataType() string {
	return "int"
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Is())
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	var val bool
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*b = NewBool(val)
	return nil
}

func NewBool(v bool) Bool {
	return Bool(v)
}
