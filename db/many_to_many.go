package db

import (
	"context"

	"github.com/baiy/gbf/log"
)

// ManyToManyBase 多对多关联模型基础字段
//
// ```sql
// CREATE TABLE `relate`
// (
//
//	`id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//	`entity_id`  bigint unsigned NOT NULL COMMENT 'config id',
//	`target_id`  bigint unsigned NOT NULL COMMENT 'auth id',
//	`created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//	PRIMARY KEY (`id`),
//	UNIQUE KEY `relate_entity_id_target_id` (`entity_id`,`target_id`),
//	KEY          `relate_target_id` (`target_id`),
//	CONSTRAINT `foreign_relate_entity_id` FOREIGN KEY (`entity_id`) REFERENCES `entity_table` (`id`) ON DELETE CASCADE,
//	CONSTRAINT `foreign_relate_target_id` FOREIGN KEY (`target_id`) REFERENCES `target_table` (`id`) ON DELETE CASCADE
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='entity - relate table';
// ```
type ManyToManyBase struct {
	Primary
	Created
	EntityId uint `json:"entity_id"`
	TargetId uint `json:"target_id"`
}

type ManyToMany[R any, E any, T any] struct {
	logger log.Logger
	db     *Db
	ctx    context.Context
}

func (m *ManyToMany[R, E, T]) Relate() []ManyToManyBase {
	return m.db.Model[ManyToManyBase](m.ctx).Table(m.db.TableName[R]()).All()
}

func (m *ManyToMany[R, E, T]) EntityIds(targetIds ...uint) []uint {
	return m.db.Model[R](m.ctx).Where("target_id in ?", targetIds).Pluck[uint]("entity_id")
}

func (m *ManyToMany[R, E, T]) Entity(targetIds ...uint) []E {
	return m.db.Model[E](m.ctx).
		Where("id in (select entity_id from %s where target_id in ? )", m.db.TableName[R](), targetIds).
		Find()
}

func (m *ManyToMany[R, E, T]) EntityPluck[F any](field string, targetIds ...uint) []F {
	return m.db.Model[E](m.ctx).
		Where("id in (select entity_id from %s where target_id in ? )", m.db.TableName[R](), targetIds).
		Pluck[F](field)
}

func (m *ManyToMany[R, E, T]) EntityLength(targetIds ...uint) int64 {
	return m.db.Model[R](m.ctx).Where("target_id in ?", targetIds).Count()
}

func (m *ManyToMany[R, E, T]) TargetIds(entityIds ...uint) []uint {
	return m.db.Model[R](m.ctx).Where("entity_id in ?", entityIds).Pluck[uint]("target_id")
}

func (m *ManyToMany[R, E, T]) Target(entityIds ...uint) []T {
	return m.db.Model[T](m.ctx).
		Where("id in (select target_id from %s where entity_id in ? )", m.db.TableName[R](), entityIds).
		Find()
}

func (m *ManyToMany[R, E, T]) TargetPluck[F any](field string, entityIds ...uint) []F {
	return m.db.Model[T](m.ctx).
		Where("id in (select target_id from "+m.db.TableName[R]()+" where entity_id in ? )", entityIds).
		Pluck[F](field)
}

func (m *ManyToMany[R, E, T]) TargetLength(entityIds ...uint) int64 {
	return m.db.Model[R](m.ctx).Where("entity_id in ?", entityIds).Count()
}

// Save 保存
func (m *ManyToMany[R, E, T]) Save(entityId, targetId uint) error {
	is := m.db.Model[R](m.ctx).Where("entity_id = ? AND target_id = ?", entityId, targetId).Exist()
	if is {
		return nil
	}
	return m.db.InsertByMap[R](m.ctx, map[string]any{"entity_id": entityId, "target_id": targetId})
}

func (m *ManyToMany[R, E, T]) SaveByTarget(targetId uint, entityIds []uint) error {
	return m.db.Transaction(func(d *Db) error {
		exist := m.EntityIds(targetId)
		add := m.diff(entityIds, exist)
		remove := m.diff(exist, entityIds)

		if len(remove) > 0 {
			d.Model[R](m.ctx).Where("entity_id in ? AND target_id = ?", remove, targetId).Delete()
		}
		if len(add) < 1 {
			return nil
		}
		var values []map[string]any
		for _, entityId := range add {
			values = append(values, map[string]any{"entity_id": entityId, "target_id": targetId})
		}
		return m.db.InsertByMap[R](m.ctx, values...)
	})
}
func (m *ManyToMany[R, E, T]) SaveByEntity(entityId uint, targetIds []uint) error {
	return m.db.Transaction(func(d *Db) error {
		exist := m.TargetIds(entityId)
		add := m.diff(targetIds, exist)
		remove := m.diff(exist, targetIds)
		if len(remove) > 0 {
			d.Model[R](m.ctx).Where("entity_id = ? AND target_id in ?", entityId, remove).Delete()
		}
		if len(add) < 1 {
			return nil
		}
		var values []map[string]any
		for _, targetId := range add {
			values = append(values, map[string]any{"entity_id": entityId, "target_id": targetId})
		}
		return m.db.InsertByMap[R](m.ctx, values...)
	})
}

func (m *ManyToMany[R, E, T]) diff(slice, comparedSlice []uint) []uint {
	var result []uint

	if len(slice) == 0 {
		return result
	}

	comparedMap := make(map[uint]struct{}, len(comparedSlice))
	for _, v := range comparedSlice {
		comparedMap[v] = struct{}{}
	}

	for _, v := range slice {
		if _, found := comparedMap[v]; !found {
			result = append(result, v)
		}
	}
	return result
}
