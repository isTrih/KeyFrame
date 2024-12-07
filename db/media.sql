CREATE TABLE `media`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `article_id` bigint(20) UNSIGNED NOT NULL,
    `cover_url`  varchar(256)        NOT NULL,
    `height`     int(11)             NOT NULL DEFAULT '0',
    `width`      int(11)             NOT NULL DEFAULT '0',
    `MediaList` JSON,
    PRIMARY KEY (`id`),
    KEY `image_post_id_645c6566_fk_article_id` (`article_id`),
    CONSTRAINT `image_post_id_645c6566_fk_article_id` FOREIGN KEY (`article_id`) REFERENCES `article` (`id`) On DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


insert into media(article_id, cover_url, height, width, MediaList)
values (1, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (2, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (3, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (4, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (5, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (6, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (7, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (8, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (9, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (10, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (11, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (12, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (13, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (14, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (15, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (16, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (17, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (18, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (19, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');
insert into media(article_id, cover_url, height, width, MediaList)
values (20, 'https://img.keaitupian.cn/newupload/05/1683621565610736.jpg', 1003, 564, '["https://img.keaitupian.cn/newupload/05/1683621565610736.jpg"]');