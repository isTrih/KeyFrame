CREATE TABLE `collect`
(
    `id`               bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`          bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
    `collected_feed_id` bigint(20) UNSIGNED NOT NULL COMMENT '被收藏帖子ID',
    `follow_status`    tinyint(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '关注状态：1-关注，2-取消关注',
    `create_time`      timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`      timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id_collected_feed_id` (`user_id`, `collected_feed_id`),
    KEY `ix_collect_user_id` (`user_id`),
    KEY `ix_collected_feed_id` (`collected_feed_id`),
    KEY `ix_update_time` (`update_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT '收藏表';