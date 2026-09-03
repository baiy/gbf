package db

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/log"
	"gorm.io/gorm"
)

// Model 通用模型
type Model[T any] struct {
	orm    gorm.ChainInterface[T]
	logger log.Logger
	db     *Db
	ctx    context.Context
}

// Where 条件
func (m *Model[T]) Where(query any, args ...any) *Model[T] {
	m.orm = m.orm.Where(query, args...)
	return m
}

// Primary 主键条件
func (m *Model[T]) Primary(value any) *Model[T] {
	m.orm = m.orm.Where("id = ?", value)
	return m
}

// Limit 查询数量
func (m *Model[T]) Limit(limit int) *Model[T] {
	m.orm = m.orm.Limit(limit)
	return m
}

// Offset 偏移
func (m *Model[T]) Offset(offset int) *Model[T] {
	m.orm = m.orm.Offset(offset)
	return m
}

// Select 查询字段
func (m *Model[T]) Select(query string, args ...any) *Model[T] {
	m.orm = m.orm.Select(query, args...)
	return m
}

// Omit 排除字段
func (m *Model[T]) Omit(columns ...string) *Model[T] {
	m.orm = m.orm.Omit(columns...)
	return m
}

// Distinct 去重
func (m *Model[T]) Distinct(args ...any) *Model[T] {
	m.orm = m.orm.Distinct(args...)
	return m
}

// Group 分组
func (m *Model[T]) Group(name string) *Model[T] {
	m.orm = m.orm.Group(name)
	return m
}

// Having 分组聚合结果过滤
func (m *Model[T]) Having(query any, args ...any) *Model[T] {
	m.orm = m.orm.Having(query, args...)
	return m
}

// Order 排序
func (m *Model[T]) Order(value any) *Model[T] {
	if v, ok := value.(OrderRequest); ok {
		value = cmp.Or(v.GetOrder(), "id") + " " + cmp.Or(v.GetDirection(), "desc")
	}
	m.orm = m.orm.Order(value)
	return m
}

// Table 指定表名
func (m *Model[T]) Table(name string, args ...any) *Model[T] {
	m.orm = m.orm.Table(name, args...)
	return m
}

// TableName 获取当前表名
func (m *Model[T]) TableName() string {
	return m.db.TableName[T]()
}

// Count 获取条数
func (m *Model[T]) Count() int64 {
	c, err := m.orm.Count(m.ctx, "*")
	if err != nil {
		panic(err)
	}
	return c
}

// Exist 是否存在
func (m *Model[T]) Exist() bool {
	return m.Count() > 0
}

// Get 获取单条
func (m *Model[T]) Get() *T {
	c, err := m.orm.Take(m.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}
	return &c
}

// GetOrNew 获取单条 不存返回 fns 参数方法数据
func (m *Model[T]) GetOrNew(fns ...func(v *T) (*T, error)) T {
	v := m.Get()
	if v == nil {
		if len(fns) == 0 {
			return *new(T)
		}
		nv, err := fns[0](new(T))
		if err != nil {
			panic(err)
		}
		return *nv
	}
	return *v
}

// GetOrFail 获取单条 不存在时插入报错
func (m *Model[T]) GetOrFail() T {
	v := m.Get()
	if v == nil {
		panic(fmt.Errorf("查询数据不存在"))
	}
	return *v
}

// Max 获取最大值
func (m *Model[T]) Max[R comparable](field string) R {
	v, _ := m.with().Value[R](fmt.Sprintf("max(%s)", field))
	return v
}

// Min 获取最小值
func (m *Model[T]) Min[R comparable](field string) R {
	v, _ := m.with().Value[R](fmt.Sprintf("min(%s)", field))
	return v
}

// Min 获取最小值
func (m *Model[T]) Sum[R comparable](field string) R {
	v, _ := m.with().Value[R](fmt.Sprintf("sum(%s)", field))
	return v
}

// Value 获取单条指定字段值
func (m *Model[T]) Value[R any](field string) (R, bool) {
	var result *R
	row := m.orm.Table(m.db.TableName[T]()).Select(field).Limit(1).Row(m.ctx)
	if err := row.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return *result, false
		} else {
			panic(err)
		}
	}
	if result == nil {
		var z R
		return z, true
	}
	return *result, true
}

// Pluck 获取多条指定字段值
func (m *Model[T]) Pluck[R any](field string) []R {
	result := []R{}
	rows, err := m.orm.Table(m.db.TableName[T]()).
		Select(field).
		Rows(m.ctx)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var value R
		if err := rows.Scan(&value); err != nil {
			panic(err)
		}
		result = append(result, value)
	}
	return result
}

// Map 获取多条记录的Map结果
func (m *Model[T]) Map[K comparable, V any](keyField, valueField string) map[K]V {
	result := map[K]V{}
	rows, err := m.orm.Table(m.db.TableName[T]()).
		Select(fmt.Sprintf("%s,%s", keyField, valueField)).
		Rows(m.ctx)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var key K
		var value V
		if err := rows.Scan(&key, &value); err != nil {
			panic(err)
		}
		result[key] = value
	}
	return result
}

// Find 获取多条
func (m *Model[T]) Find() []T {
	c, err := m.orm.Find(m.ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// All 获取所有记录
func (m *Model[T]) All() []T {
	return m.with().Limit(-1).Offset(-1).Find()
}

// FindInBatches 分批查询/处理
func (m *Model[T]) FindInBatches(limit int, fc func(data []T, batch int) error) {
	err := m.orm.FindInBatches(m.ctx, limit, fc)
	if err != nil {
		panic(err)
	}
}

// Paginate 分页
func (m *Model[T]) Paginate(paginate PaginateRequest) Paginate[T] {
	page := max(paginate.GetPage(), 1)
	pageSize := max(cmp.Or(paginate.GetPageSize(), 50), 1)

	r := Paginate[T]{
		Last:       true,
		Page:       page,
		PageSize:   pageSize,
		Total:      0,
		TotalPages: 0,
		Lists:      []T{},
	}

	count := m.with().Count()

	if count == 0 {
		return r
	}

	totalPages := uint(math.Ceil(float64(count) / float64(pageSize)))
	r.Last = page >= totalPages
	r.Total = uint64(count)
	r.TotalPages = totalPages
	r.Lists = m.Order(paginate).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find()
	return r
}

// Delete 删除
func (m *Model[T]) Delete() {
	_, err := m.orm.Delete(m.ctx)
	if err != nil {
		panic(err)
	}
}

// Update 更新(全字段覆盖)
func (m *Model[T]) Update(d T, columns ...string) error {
	return m.update(d, "")
}

// UpdateSelect 更新(选择字段)
func (m *Model[T]) UpdateSelect(d T, columns ...string) error {
	return m.update(d, "select", columns...)
}

// UpdateOmit 更新(排除字段)
func (m *Model[T]) UpdateOmit(d T, columns ...string) error {
	return m.update(d, "omit", columns...)
}

// UpdateColumn 更新单列
func (m *Model[T]) UpdateColumn(column string, value any) error {
	_, err := m.orm.Select("*").Update(m.ctx, column, value)
	if err != nil {
		m.logger.Errorf("update column error: %v", err)
		return errors.NewBusinessf("更新失败")
	}
	return nil
}

// 更新
// t 更新类型 "": 所有字段 "select": 更新字段 "omit": 排除字段
// columns 更新/排除字段
func (m *Model[T]) update(d T, t string, columns ...string) error {
	if t != "" && len(columns) < 1 {
		return errors.NewBusinessf("更新/排除字段不能为空")
	}
	var err error
	switch t {
	case "select":
		query := ""
		var otherFields []any
		if len(columns) > 0 {
			query = columns[0]
			for _, temp := range columns[1:] {
				otherFields = append(otherFields, temp)
			}
		}
		_, err = m.orm.Select(query, otherFields...).Updates(m.ctx, d)
	case "omit":
		_, err = m.orm.Omit(columns...).Updates(m.ctx, d)
	default:
		_, err = m.orm.Select("*").Updates(m.ctx, d)
	}
	if err != nil {
		m.logger.Errorf("update error: %v", err)
		return errors.NewBusinessf("更新失败")
	}
	return nil
}

func (m *Model[T]) with() *Model[T] {
	return newModel(m.ctx, m.orm.Where(""), m.logger, m.db)
}

type PaginateRequest interface {
	GetPage() uint
	GetPageSize() uint
	OrderRequest
}

type OrderRequest interface {
	GetOrder() string     // 排序字段
	GetDirection() string // 排序方式 asc desc
}

func newModel[T any](ctx context.Context, orm gorm.ChainInterface[T], logger log.Logger, db *Db) *Model[T] {
	return &Model[T]{orm: orm, logger: logger, db: db, ctx: ctx}
}
