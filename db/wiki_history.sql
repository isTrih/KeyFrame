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

 Date: 08/04/2025 09:34:34
*/


-- ----------------------------
-- Table structure for wiki_history
-- ----------------------------
DROP TABLE IF EXISTS "public"."wiki_history";
CREATE TABLE "public"."wiki_history" (
  "id" int8 NOT NULL DEFAULT nextval('wiki_id_seq'::regclass),
  "wiki_id" int8 NOT NULL,
  "raw_content" text COLLATE "pg_catalog"."default",
  "editor_id" int8 NOT NULL,
  "edit_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "change_log" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."wiki_history" OWNER TO "keyframe_core";
