package db

// 数据操作

import (
	"context"
	"strings"

	"github.com/baiy/gbf/errors"
	"gorm.io/gorm"
)

// Save 保存数据
// 存在主键更新 不存在插入
func (d *Db) Save[T any](ctx context.Context, value *T) error {
	result := d.Gorm().WithContext(ctx).Save(value)
	if result.Error != nil {
		d.logger.Errorf("db save error: %v", result.Error)
		return errors.NewBusiness("数据保存数据库失败")
	}
	return nil
}

// Insert 插入数据
func (d *Db) Insert[T any](ctx context.Context, value *T) error {
	err := gorm.G[T](d.Gorm()).Create(ctx, value)
	if err != nil {
		d.logger.Errorf("db insert error: %v", err)
		return errors.NewBusiness("数据插入数据库失败")
	}
	return nil
}

// InsertByMap 插入数据
func (d *Db) InsertByMap[T any](ctx context.Context, values ...map[string]any) error {
	result := d.Gorm().WithContext(ctx).Model(new(T)).Create(values)
	if result.Error != nil {
		d.logger.Errorf("db insertByMap error: %v", result.Error)
		return errors.NewBusiness("数据插入数据库失败")
	}
	return nil
}

// InsertMultiple 插入多条数据
func (d *Db) InsertMultiple[T any](ctx context.Context, values []T, batchSizes ...int) error {
	var err error
	if len(batchSizes) > 0 && batchSizes[0] > 0 {
		err = gorm.G[T](d.Gorm()).CreateInBatches(ctx, &values, batchSizes[0])
	} else {
		err = d.Gorm().WithContext(ctx).Create(&values).Error
	}
	if err != nil {
		d.logger.Errorf("db insertMultiple error: %v", err)
		return errors.NewBusiness("数据批量插入数据库失败")
	}
	return nil
}

// DeleteById 根据主键删除
func (d *Db) DeleteById[T any](ctx context.Context, ids ...uint) {
	if len(ids) == 0 {
		panic("ids不能为空")
	}
	d.Model[T](ctx).Where("id in ?", ids).Delete()
}

// ExecMultipleSql 执行多条sql
//
// sql 多条sql语句 使用";"分割
// tablePrefix 表前缀
func (d *Db) ExecMultipleSql(sql, tablePrefix string) error {
	return d.Gorm().Connection(func(tx *gorm.DB) error {
		if tablePrefix != "" {
			sql = strings.ReplaceAll(sql, tablePrefix, d.TableName(""))
		}
		sqlStatements := strings.Split(sql, ";")
		tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		for _, v := range sqlStatements {
			sql := strings.TrimSpace(v)
			if sql == "" {
				continue
			}
			result := tx.Exec(sql)
			if result.Error != nil {
				return result.Error
			}
		}
		return tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	})
}

// Exec 执行sql
func (d *Db) Exec(sql, tablePrefix string) error {
	if tablePrefix != "" {
		sql = strings.ReplaceAll(sql, tablePrefix, d.TableName(""))
	}
	return d.Gorm().Exec(sql).Error
}
