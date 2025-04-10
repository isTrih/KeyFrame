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

 Date: 10/04/2025 17:14:39
*/


-- ----------------------------
-- Table structure for comment_metrics
-- ----------------------------
DROP TABLE IF EXISTS "public"."comment_metrics";
CREATE TABLE "public"."comment_metrics" (
  "id" int8 NOT NULL DEFAULT nextval('article_metrics_id_seq'::regclass),
  "comment_id" int8 NOT NULL DEFAULT 0,
  "likes" int8 NOT NULL DEFAULT 0,
  "comments" int8 NOT NULL DEFAULT 0
)
;
ALTER TABLE "public"."comment_metrics" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."comment_metrics"."comment_id" IS '评论ID';
COMMENT ON COLUMN "public"."comment_metrics"."likes" IS '点赞数';
COMMENT ON COLUMN "public"."comment_metrics"."comments" IS '评论数';

-- ----------------------------
-- Indexes structure for table comment_metrics
-- ----------------------------
CREATE UNIQUE INDEX "idx_comment_metrics_comment_id" ON "public"."comment_metrics" USING btree (
  "comment_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_comment_metrics_likes" ON "public"."comment_metrics" USING btree (
  "likes" "pg_catalog"."int8_ops" DESC NULLS LAST
);
CREATE INDEX "uk_comment_id_copy1" ON "public"."comment_metrics" USING btree (
  "comment_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table comment_metrics
-- ----------------------------
ALTER TABLE "public"."comment_metrics" ADD CONSTRAINT "article_metrics_copy1_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table comment_metrics
-- ----------------------------
ALTER TABLE "public"."comment_metrics" ADD CONSTRAINT "article_metrics_copy1_article_id_fkey" FOREIGN KEY ("comment_id") REFERENCES "public"."article" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
