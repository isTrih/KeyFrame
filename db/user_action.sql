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

 Date: 13/04/2025 01:02:30
*/


-- ----------------------------
-- Table structure for user_action
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_action";
CREATE TABLE "public"."user_action" (
  "id" int8 NOT NULL DEFAULT nextval('user_action_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "target_id" int8 NOT NULL,
  "action_type" int2 NOT NULL,
  "target_type" int2 NOT NULL,
  "action_value" int2 NOT NULL DEFAULT 1,
  "create_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."user_action" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."user_action"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."user_action"."target_id" IS '操作目标ID（文章/评论/用户等）';
COMMENT ON COLUMN "public"."user_action"."action_type" IS '操作类型：1-动态点赞, 2-评论点赞, 3-收藏, 4-关注';
COMMENT ON COLUMN "public"."user_action"."target_type" IS '目标类型：1-文章, 2-评论, 3-用户';
COMMENT ON COLUMN "public"."user_action"."action_value" IS '操作值（1-执行，0-取消）';
COMMENT ON TABLE "public"."user_action" IS '用户操作记录表';

-- ----------------------------
-- Indexes structure for table user_action
-- ----------------------------
CREATE INDEX "idx_target_actions" ON "public"."user_action" USING btree (
  "target_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "action_type" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "target_type" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_actions" ON "public"."user_action" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "action_type" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "target_type" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "create_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_user_action" ON "public"."user_action" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "target_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "action_type" "pg_catalog"."int2_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table user_action
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."user_action"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table user_action
-- ----------------------------
ALTER TABLE "public"."user_action" ADD CONSTRAINT "user_action_pkey" PRIMARY KEY ("id");
