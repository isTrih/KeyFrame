package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"net"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		UpdateIp(ctx context.Context, id uint64, il string, ia string) error
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

func (m *customUserModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// UpdateIp 更新用户ip
//
//		@id: 用户ID
//		@ia: 用户IP地址
//	 @il: 用户IP归属地
//
// returns error
func (m *defaultUserModel) UpdateIp(ctx context.Context, id uint64, il string, ia string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	ipAddressBinary := net.ParseIP(ia)
	chaozjUserIdKey := fmt.Sprintf("%s%v", cacheChaozjUserIdPrefix, data.Id)
	chaozjUserMobileKey := fmt.Sprintf("%s%v", cacheChaozjUserMobilePrefix, data.Mobile)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {

		query := fmt.Sprintf("update %s set `ip_location` = ?, `ip_address` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, il, ipAddressBinary, id)
	}, chaozjUserIdKey, chaozjUserMobileKey)
	return err
}
