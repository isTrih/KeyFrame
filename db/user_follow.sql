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

 Date: 13/04/2025 23:25:46
*/


-- ----------------------------
-- Table structure for user_follow
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_follow";
CREATE TABLE "public"."user_follow" (
  "id" int8 NOT NULL DEFAULT nextval('user_follow_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "followed_user_id" int8 NOT NULL,
  "status" int2 NOT NULL DEFAULT 1,
  "create_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."user_follow" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."user_follow"."user_id" IS '关注者ID';
COMMENT ON COLUMN "public"."user_follow"."followed_user_id" IS '被关注用户ID';
COMMENT ON COLUMN "public"."user_follow"."status" IS '状态：1-关注，0-取消';
COMMENT ON TABLE "public"."user_follow" IS '用户关注关系表';

-- ----------------------------
-- Indexes structure for table user_follow
-- ----------------------------
CREATE INDEX "idx_user_follow_followed_status" ON "public"."user_follow" USING btree (
  "followed_user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_follow_user_status" ON "public"."user_follow" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_followers" ON "public"."user_follow" USING btree (
  "followed_user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "create_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "idx_user_follows" ON "public"."user_follow" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "create_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_follow_relation" ON "public"."user_follow" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "followed_user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table user_follow
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."user_follow"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table user_follow
-- ----------------------------
ALTER TABLE "public"."user_follow" ADD CONSTRAINT "user_follow_pkey" PRIMARY KEY ("id");
