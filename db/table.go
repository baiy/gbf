package db

import "gorm.io/gorm"

// TableName 获取表名
// name 传值时字符串转表名 否则按照泛型T model转表名
func (d *Db) TableName[T any](names ...T) string {
	// 字符串转表名
	if len(names) == 1 {
		if name, ok := any(names[0]).(string); ok {
			return d.Gorm().NamingStrategy.TableName(name)
		}
		panic("TableName: invalid name type")
		// 其他值 按照model转表名
	}

	// model转表名
	stmt := &gorm.Statement{DB: d.Gorm()}
	err := stmt.Parse(new(T))
	if err != nil {
		d.logger.Errorf("parse table error:%v", err)
	}
	return stmt.Schema.Table
}

// HasTable 判断表是否存在
func (d *Db) HasTable[T any](names ...T) bool {
	return d.Gorm().Migrator().HasTable(d.TableName[T](names...))
}

// DropTable 删除表
func (d *Db) DropTable[T any](names ...T) error {
	return d.Gorm().Migrator().DropTable(d.TableName[T](names...))
}
