CREATE TABLE `user`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `password`    varchar(256)        NOT NULL COMMENT '密码',
    `nickname`    varchar(32)         NOT NULL DEFAULT '' COMMENT '昵称',
    `signature`   varchar(128)        NOT NULL DEFAULT '该用户是一个慵懒动画人，还未设置简介' COMMENT '简介',
    `avatar`      varchar(256)        NOT NULL DEFAULT '' COMMENT '头像',
    `type`        tinyint(4)          NOT NULL DEFAULT '0' COMMENT '状态 0:默认用户 1:正式用户 2:V认证用户 3:管理员用户',
    `vnote`       varchar(256)        NOT NULL DEFAULT '' COMMENT 'V认证信息',
    `mobile`      varchar(128)        NOT NULL DEFAULT '' COMMENT '手机号',
    `status`      tinyint(4)          NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:禁用 2:删除',
    `banned_time` timestamp COMMENT '封禁时间',
    `ban_time`    int(0)             NOT NULL DEFAULT '0' COMMENT '封禁时长',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_update_time` (`update_time`),
    UNIQUE KEY `uk_mobile` (`mobile`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户表';

insert into user(nickname, password, avatar, mobile)
values ('三氢超正经', '96cae35ce8a9b0244178bf28e4966c2ce1b8385723a96a6b838858cdd6ca0a1e',
        'https://coss.chaozj.com/avatar/test_1.jpg', '13800138000');
INSERT INTO `chaozj`.`user` (`id`, `password`, `nickname`, `signature`, `avatar`, `type`, `mobile`, `create_time`,
                             `update_time`)
VALUES (10, 'd6ab66d05bb8659b4794c29d75dedd147105bb80ad312f80f08820ba9510c417', '楚盒汉界', 'CHAOZJ',
        'https://coss.chaozj.com/avatar/test_4.jpg', 0, '18655311015', '2024-11-26 13:01:39',
        '2024-12-02 19:21:18');
INSERT INTO user (`id`, `password`, `nickname`, `signature`, `avatar`, `type`, `mobile`, `create_time`,
                             `update_time`)
VALUES (11, '364f2f9ffa3564fd0750be9a06ddd7878b7385f67328c73a17e9fa646d93c6c0', '龙虾秘书', 'CHAOZJ',
        'https://coss.chaozj.com/avatar/test_3.jpg', 0, '19885452781', '2024-11-27 23:32:51',
        '2024-12-02 19:21:22');

INSERT INTO user (`id`, `password`, `nickname`, `signature`, `avatar`, `type`, `mobile`, `create_time`,
                             `update_time`)
VALUES (88, '364f2f9ffa3564fd0750be9a06ddd7878b7385f67328c73a17e9fa646d93c6c0', '魔法飞鱼', 'CHAOZJ',
        'https://coss.chaozj.com/avatar/test_2.jpg', 0, '19885452781', '2024-11-27 23:32:51',
        '2024-12-02 19:21:22');