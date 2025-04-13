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

 Date: 13/04/2025 01:09:10
*/


-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS "public"."article";
CREATE TABLE "public"."article" (
  "id" int8 NOT NULL DEFAULT nextval('article_id_seq'::regclass),
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "raw_content" text COLLATE "pg_catalog"."default" NOT NULL,
  "author_id" int8 NOT NULL,
  "status" int2 NOT NULL DEFAULT 0,
  "type" int2 NOT NULL DEFAULT 0,
  "ip_location" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "publish_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "create_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "ai_insp" int2 NOT NULL DEFAULT 0,
  "insp" int2 NOT NULL DEFAULT 1
)
;
ALTER TABLE "public"."article" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."article"."id" IS '主键ID';
COMMENT ON COLUMN "public"."article"."title" IS '标题';
COMMENT ON COLUMN "public"."article"."content" IS '原始文本';
COMMENT ON COLUMN "public"."article"."raw_content" IS '内容';
COMMENT ON COLUMN "public"."article"."author_id" IS '作者ID';
COMMENT ON COLUMN "public"."article"."status" IS '状态 0:正常 1:仅自己 2:删除';
COMMENT ON COLUMN "public"."article"."type" IS '状态 0:默认图片帧 1:视频帧 2:纯文字帧';
COMMENT ON COLUMN "public"."article"."ip_location" IS 'IP归属地';
COMMENT ON COLUMN "public"."article"."publish_time" IS '发布时间';
COMMENT ON COLUMN "public"."article"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."article"."update_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."article"."ai_insp" IS 'AI检验，默认0没有问题，不为0出现问题。';
COMMENT ON COLUMN "public"."article"."insp" IS '人工校验为0时可用，默认为1待检验';
COMMENT ON TABLE "public"."article" IS '动态表';

-- ----------------------------
-- Indexes structure for table article
-- ----------------------------
CREATE INDEX "idx_article_raw_content_gin" ON "public"."article" USING gin (
  to_tsvector('zh_cn'::regconfig, raw_content) "pg_catalog"."tsvector_ops"
);
CREATE INDEX "idx_article_status_type" ON "public"."article" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "type" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE INDEX "idx_article_update_time" ON "public"."article" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" DESC NULLS LAST
);
CREATE INDEX "ix_author_id" ON "public"."article" USING btree (
  "author_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "ix_update_time" ON "public"."article" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table article
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."article"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table article
-- ----------------------------
ALTER TABLE "public"."article" ADD CONSTRAINT "article_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table article
-- ----------------------------
ALTER TABLE "public"."article" ADD CONSTRAINT "ix_author_id" FOREIGN KEY ("author_id") REFERENCES "public"."user" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
