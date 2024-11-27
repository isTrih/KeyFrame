CREATE TABLE `user`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `password`    varchar(256)        NOT NULL COMMENT '密码',
    `nickname`    varchar(32)         NOT NULL DEFAULT '' COMMENT '昵称',
    `signature`   varchar(128)        NOT NULL DEFAULT '该用户是一个慵懒动画人，还未设置简介' COMMENT '简介',
    `avatar`      varchar(256)        NOT NULL DEFAULT '' COMMENT '头像',
    `type`        tinyint(4)          NOT NULL DEFAULT '0' COMMENT '状态 0:默认用户 1:正式用户 2:V认证用户 3:管理员用户',
    `mobile`      varchar(128)        NOT NULL DEFAULT '' COMMENT '手机号',
    `create_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_update_time` (`update_time`),
    UNIQUE KEY `uk_mobile` (`mobile`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='用户表';

insert into user(nickname, password, avatar, mobile)
values ('张三','96cae35ce8a9b0244178bf28e4966c2ce1b8385723a96a6b838858cdd6ca0a1e', 'https://beyond-blog.oss-cn-beijing.aliyuncs.com/avatar/2021/01/01/1609488000.jpg', '13800138000');