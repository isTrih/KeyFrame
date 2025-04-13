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

 Date: 14/04/2025 02:24:13
*/


-- ----------------------------
-- Table structure for wiki
-- ----------------------------
DROP TABLE IF EXISTS "public"."wiki";
CREATE TABLE "public"."wiki" (
  "id" int8 NOT NULL DEFAULT nextval('wiki_id_seq'::regclass),
  "type" varchar COLLATE "pg_catalog"."default",
  "raw_content" text COLLATE "pg_catalog"."default",
  "content" text COLLATE "pg_catalog"."default",
  "editor" json,
  "status" int2,
  "name" json,
  "create_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."wiki" OWNER TO "keyframe_core";

-- ----------------------------
-- Triggers structure for table wiki
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."wiki"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table wiki
-- ----------------------------
ALTER TABLE "public"."wiki" ADD CONSTRAINT "wiki_pkey" PRIMARY KEY ("id");
