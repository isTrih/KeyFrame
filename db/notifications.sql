/*
 Navicat Premium Data Transfer

 Source Server         : PG
 Source Server Type    : PostgreSQL
 Source Server Version : 160008 (160008)
 Source Host           : 121.37.154.241:5432
 Source Catalog        : keyframe_test
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160008 (160008)
 File Encoding         : 65001

 Date: 14/04/2025 06:55:19
*/


-- ----------------------------
-- Table structure for notifications
-- ----------------------------
DROP TABLE IF EXISTS "public"."notifications";
CREATE TABLE "public"."notifications" (
  "id" int8 NOT NULL DEFAULT nextval('notifications_id_seq'::regclass),
  "sender_id" int8,
  "receiver_id" int8 NOT NULL,
  "type" int2 NOT NULL,
  "content" text COLLATE "pg_catalog"."default",
  "target_id" int8,
  "is_read" bool DEFAULT false,
  "create_time" timestamptz(6) DEFAULT CURRENT_TIMESTAMP,
  "extra" jsonb DEFAULT '{}'::jsonb,
  "target_content" text COLLATE "pg_catalog"."default",
  "target_type" int2 NOT NULL
)
;
ALTER TABLE "public"."notifications" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."notifications"."id" IS '主键，自增 ID';
COMMENT ON COLUMN "public"."notifications"."sender_id" IS '发送人 ID，系统通知可为 NULL';
COMMENT ON COLUMN "public"."notifications"."receiver_id" IS '接收人 ID';
COMMENT ON COLUMN "public"."notifications"."type" IS '通知类型，例如 1点赞2评论3收藏4关注';
COMMENT ON COLUMN "public"."notifications"."content" IS '通知显示内容，可用于消息列表展示';
COMMENT ON COLUMN "public"."notifications"."target_id" IS '通知对象的 ID';
COMMENT ON COLUMN "public"."notifications"."is_read" IS '是否已读标志';
COMMENT ON COLUMN "public"."notifications"."create_time" IS '通知创建时间';
COMMENT ON COLUMN "public"."notifications"."extra" IS '额外信息（JSON 格式），如帖子标题、评论摘要等';
COMMENT ON COLUMN "public"."notifications"."target_content" IS '目标内容';
COMMENT ON COLUMN "public"."notifications"."target_type" IS '目标类型 文章，评论，用户';
COMMENT ON TABLE "public"."notifications" IS '用户通知表，记录系统消息和用户互动通知';

-- ----------------------------
-- Indexes structure for table notifications
-- ----------------------------
CREATE INDEX "idx_notifications_is_read" ON "public"."notifications" USING btree (
  "receiver_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "is_read" "pg_catalog"."bool_ops" ASC NULLS LAST
);
CREATE INDEX "idx_notifications_receiver_id" ON "public"."notifications" USING btree (
  "receiver_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table notifications
-- ----------------------------
ALTER TABLE "public"."notifications" ADD CONSTRAINT "notifications_pkey" PRIMARY KEY ("id");
