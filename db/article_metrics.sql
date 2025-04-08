/*
 Navicat Premium Data Transfer

 Source Server         : PG
 Source Server Type    : PostgreSQL
 Source Server Version : 160006 (160006)
 Source Host           : 121.37.154.241:5432
 Source Catalog        : keyframe_test
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160006 (160006)
 File Encoding         : 65001

 Date: 08/04/2025 08:18:00
*/


-- ----------------------------
-- Table structure for article_metrics
-- ----------------------------
DROP TABLE IF EXISTS "public"."article_metrics";
CREATE TABLE "public"."article_metrics" (
  "id" int8 NOT NULL DEFAULT nextval('article_metrics_id_seq'::regclass),
  "article_id" int8 NOT NULL DEFAULT 0,
  "likes" int8 NOT NULL DEFAULT 0,
  "favorites" int8 NOT NULL DEFAULT 0,
  "comments" int8 NOT NULL DEFAULT 0,
  "shares" int8 NOT NULL DEFAULT 0
)
;
ALTER TABLE "public"."article_metrics" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."article_metrics"."article_id" IS '文章ID';
COMMENT ON COLUMN "public"."article_metrics"."likes" IS '点赞数';
COMMENT ON COLUMN "public"."article_metrics"."favorites" IS '收藏数';
COMMENT ON COLUMN "public"."article_metrics"."comments" IS '评论数';
COMMENT ON COLUMN "public"."article_metrics"."shares" IS '分享数';

-- ----------------------------
-- Indexes structure for table article_metrics
-- ----------------------------
CREATE INDEX "uk_article_id" ON "public"."article_metrics" USING btree (
  "article_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table article_metrics
-- ----------------------------
ALTER TABLE "public"."article_metrics" ADD CONSTRAINT "article_metrics_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table article_metrics
-- ----------------------------
ALTER TABLE "public"."article_metrics" ADD CONSTRAINT "article_metrics_ibfk_1" FOREIGN KEY ("article_id") REFERENCES "public"."article" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
