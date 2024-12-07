CREATE TABLE `article`
(
    `id`           bigint(20) UNSIGNED          NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `title`        varchar(255)                 NOT NULL DEFAULT '' COMMENT '标题',
    `content`      text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
    `author_id`    bigint(20) UNSIGNED          NOT NULL DEFAULT '0' COMMENT '作者ID',
    `status`       tinyint(4)                   NOT NULL DEFAULT '0' COMMENT '状态 0:待审核 1:审核不通过 2:可见 3:用户隐私',
    `type`         tinyint(4)                   NOT NULL DEFAULT '0' COMMENT '状态 0:默认图片帧 1:视频帧 2:纯文字帧',
    `comment_num`  int(11)                      NOT NULL DEFAULT '0' COMMENT '评论数',
    `like_num`     int(11)                      NOT NULL DEFAULT '0' COMMENT '点赞数',
    `collect_num`  int(11)                      NOT NULL DEFAULT '0' COMMENT '收藏数',
    `view_num`     int(11)                      NOT NULL DEFAULT '0' COMMENT '浏览数',
    `share_num`    int(11)                      NOT NULL DEFAULT '0' COMMENT '分享数',
    `tag_ids`      varchar(255)                 NOT NULL DEFAULT '' COMMENT '标签ID',
    `publish_time` timestamp                    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
    `create_time`  timestamp                    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  timestamp                    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_author_id` (`author_id`),
    KEY `ix_update_time` (`update_time`),
    CONSTRAINT `ix_author_id` FOREIGN KEY (`author_id`) REFERENCES `user` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT ='动态表';


insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题1', '文章内容1', 1, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题2', '文章内容2', 1, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题3', '文章内容3', 10, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题4', '文章内容4', 10, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题5', '文章内容5', 10, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题6', '文章内容6', 10, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题7', '文章内容7', 11, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题8', '文章内容8', 11, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题9', '文章内容9', 11, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题10', '文章内容10', 11, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题11', '文章内容11', 1, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题12', '文章内容12', 1, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题13', '文章内容13', 10, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题14', '文章内容14', 10, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题15', '文章内容15', 10, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题16', '文章内容16', 10, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题17', '文章内容17', 11, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题18', '文章内容18', 11, 10, '2023-11-25 15:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题19', '文章内容19', 11, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time)
values ('文章标题20', '文章内容20', 11, 10, '2023-11-25 15:01:01');



