CREATE TABLE `prefix_cache`
(

    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(150) NOT NULL COMMENT '缓存标识',
    `value`      mediumtext   NOT NULL COMMENT '缓存值',
    `expire`     datetime COMMENT '过期时间',
    `created_at` datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据缓存表';
