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

 Date: 10/04/2025 17:14:22
*/


-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS "public"."comment";
CREATE TABLE "public"."comment" (
  "id" int8 NOT NULL DEFAULT nextval('comment_id_seq'::regclass),
  "article_id" int8 NOT NULL,
  "parent_id" int8 NOT NULL DEFAULT 0,
  "parent_user_id" int8 NOT NULL DEFAULT 0,
  "user_id" int8 NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 NOT NULL DEFAULT 0,
  "insp" int2 NOT NULL DEFAULT 1,
  "ai_insp" int2 NOT NULL DEFAULT 0,
  "ip_location" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "create_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."comment" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."comment"."id" IS '主键ID';
COMMENT ON COLUMN "public"."comment"."article_id" IS '归属文章ID';
COMMENT ON COLUMN "public"."comment"."parent_id" IS '评论目标id，为0则回复文章';
COMMENT ON COLUMN "public"."comment"."parent_user_id" IS '被回复用户ID，为-1则无被回复id';
COMMENT ON COLUMN "public"."comment"."user_id" IS '评论用户ID';
COMMENT ON COLUMN "public"."comment"."content" IS '内容';
COMMENT ON COLUMN "public"."comment"."status" IS '状态 0:正常 1:删除';
COMMENT ON COLUMN "public"."comment"."insp" IS '人工校验为0时可用，默认为1待检验';
COMMENT ON COLUMN "public"."comment"."ai_insp" IS 'AI检验，默认0没有问题，1出现问题。';
COMMENT ON COLUMN "public"."comment"."ip_location" IS 'IP归属地';
COMMENT ON COLUMN "public"."comment"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."comment"."update_time" IS '最后修改时间';
COMMENT ON TABLE "public"."comment" IS '评论表';

-- ----------------------------
-- Indexes structure for table comment
-- ----------------------------
CREATE INDEX "fk_comment_article_id" ON "public"."comment" USING btree (
  "article_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_comment_update_time" ON "public"."comment" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" DESC NULLS LAST
);
CREATE INDEX "ix_update_time_copy_3" ON "public"."comment" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table comment
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."comment"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table comment
-- ----------------------------
ALTER TABLE "public"."comment" ADD CONSTRAINT "comment_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table comment
-- ----------------------------
ALTER TABLE "public"."comment" ADD CONSTRAINT "fk_comment_article_id" FOREIGN KEY ("article_id") REFERENCES "public"."article" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
