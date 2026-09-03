package db

// 模型常见字段
// CREATE TABLE `name`
// (
//    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
//    `created_at`  datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//    `updated_at`  datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
//    PRIMARY KEY (`id`),
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='table comment';

import (
	"github.com/baiy/gbf/db/mysql/field"
)

type Base struct {
	Primary
	Created
	Updated
}

type Primary struct {
	ID uint `gorm:"primaryKey" json:"id"`
}

func (p Primary) GetId() uint {
	return p.ID
}

type Created struct {
	// 不自动更新 由数据库维护
	CreatedAt field.DateTimeMillisecond `gorm:"->" json:"createdAt"`
}

type Updated struct {
	// 不自动更新 由数据库维护
	UpdatedAt field.DateTimeMillisecond `gorm:"->" json:"updatedAt"`
}
