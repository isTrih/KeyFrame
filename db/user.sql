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

 Date: 10/04/2025 17:13:59
*/


-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "public"."user";
CREATE TABLE "public"."user" (
  "id" int8 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
  "password" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "nickname" varchar(12) COLLATE "pg_catalog"."default" NOT NULL,
  "signature" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'CHAOZJ'::character varying,
  "avatar" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "type" int4 NOT NULL DEFAULT 0,
  "vnote" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "mobile" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 NOT NULL DEFAULT 0,
  "banned_time" timestamp(6),
  "ban_time" int4 NOT NULL DEFAULT 0,
  "create_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "ip_address" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ip_location" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
ALTER TABLE "public"."user" OWNER TO "keyframe_core";
COMMENT ON COLUMN "public"."user"."id" IS '主键ID';
COMMENT ON COLUMN "public"."user"."password" IS '密码';
COMMENT ON COLUMN "public"."user"."nickname" IS '昵称';
COMMENT ON COLUMN "public"."user"."signature" IS '简介';
COMMENT ON COLUMN "public"."user"."avatar" IS '头像';
COMMENT ON COLUMN "public"."user"."type" IS ' 0:默认用户 100:正式用户 200:个人V认证用户 300:企业V认证用户，400以上:超正经员工, 900以上特殊账号';
COMMENT ON COLUMN "public"."user"."vnote" IS 'V认证信息';
COMMENT ON COLUMN "public"."user"."mobile" IS '手机号';
COMMENT ON COLUMN "public"."user"."status" IS '状态 0:正常 1:禁用 2:删除';
COMMENT ON COLUMN "public"."user"."banned_time" IS '封禁时间';
COMMENT ON COLUMN "public"."user"."ban_time" IS '封禁时长';
COMMENT ON COLUMN "public"."user"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."user"."update_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."user"."ip_address" IS 'IP地址';
COMMENT ON COLUMN "public"."user"."ip_location" IS 'IP归属地';
COMMENT ON TABLE "public"."user" IS '用户表';

-- ----------------------------
-- Indexes structure for table user
-- ----------------------------
CREATE INDEX "idx_user_mobile" ON "public"."user" USING btree (
  "mobile" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "ix_update_time_copy_1" ON "public"."user" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_mobile" ON "public"."user" USING btree (
  "mobile" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table user
-- ----------------------------
CREATE TRIGGER "update_timestamp_trigger" BEFORE INSERT OR UPDATE ON "public"."user"
FOR EACH ROW
EXECUTE PROCEDURE "public"."update_timestamp"();

-- ----------------------------
-- Primary Key structure for table user
-- ----------------------------
ALTER TABLE "public"."user" ADD CONSTRAINT "user_pkey" PRIMARY KEY ("id");
