package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UpdateIp(ctx context.Context, id uint64, il string, ia string) error
		UpdateIpByMobile(ctx context.Context, mobile string, il string, ia string) error
		UpdatePassword(ctx context.Context, id uint64, newPassword string) error
		CreateUserByCZJCode(ctx context.Context, czjCode int64, name, password, mobile string) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

// UpdateIp 更新用户ip
//
//		@id: 用户ID
//		@ia: 用户IP地址
//	 @il: 用户IP归属地
//
// returns error
func (m *defaultUserModel) UpdateIp(ctx context.Context, id uint64, il string, ia string) error {

	chaozjUserIdKey := fmt.Sprintf("%s%v", publicUserIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(`update %s set ip_location = $1, ip_address =$2 where id = $3`, m.table)
		return conn.ExecCtx(ctx, query, il, ia, id)
	}, chaozjUserIdKey)
	return err
}

// UpdateIpByMobile 手机号用更新用户ip
//
//		@mobile: 用户手机号
//		@ia: 用户IP地址
//	 @il: 用户IP归属地
//
// returns error
func (m *defaultUserModel) UpdateIpByMobile(ctx context.Context, mobile string, il string, ia string) error {
	chaozjUserMobileKey := fmt.Sprintf("%s%v", publicUserMobilePrefix, mobile)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(`update %s set ip_location = $1, ip_address =$2 where mobile = $3`, m.table)
		return conn.ExecCtx(ctx, query, il, ia, mobile)
	}, chaozjUserMobileKey)
	return err
}

func (m *defaultUserModel) UpdatePassword(ctx context.Context, id uint64, newPassword string) error {
	chaozjUserIdKey := fmt.Sprintf("%s%v", publicUserIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE %s SET `password` = ? WHERE `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, newPassword, id)
	}, chaozjUserIdKey)
	return err
}

func (m *defaultUserModel) CreateUserByCZJCode(ctx context.Context, czjCode int64, name, password, mobile string) error {
	query := `
	INSERT INTO "public"."user" (
    "id",
    "password",
    "nickname",
    "signature",
    "avatar",
    "type",
    "vnote",
    "mobile",
    "status",
    "banned_time",
    "ban_time",
    "create_time",
    "update_time",
    "ip_address",
    "ip_location"
) VALUES (
    $1,  -- 强制使用的ID
    $2,  -- 加密后的密码
    $3,  -- 用户名
    '',  -- 默认签名
    'avatar.jpg',  -- 默认头像
    0,  -- 默认用户类型
    '',  -- 空V认证信息
    $4,  -- 手机号
    0,  -- 默认状态(正常)
    NULL,  -- 无封禁时间
    0,  -- 默认封禁时长
    CURRENT_TIMESTAMP,  -- 当前时间
    CURRENT_TIMESTAMP,  -- 当前时间
    '',  -- 空IP地址
    ''  -- 空IP归属地
)
ON CONFLICT (id) DO NOTHING;
	`
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, czjCode, password, name, mobile)
	})
	if err != nil {
		return err
	}
	return nil
}
