package user_follow

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserFollowModel = (*customUserFollowModel)(nil)

type (
	// UserFollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserFollowModel.
	UserFollowModel interface {
		userFollowModel
		GetUserFollowNum(ctx context.Context, uid uint64) (*UserFollowNum, error)
	}

	customUserFollowModel struct {
		*defaultUserFollowModel
	}

	UserFollowNum struct {
		FollowCount int `db:"follow_count"`   // 关注
		FanCount    int `db:"follower_count"` // 粉丝
	}
)

// NewUserFollowModel returns a model for the database table.
func NewUserFollowModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserFollowModel {
	return &customUserFollowModel{
		defaultUserFollowModel: newUserFollowModel(conn, c, opts...),
	}
}

func (m *defaultUserFollowModel) GetUserFollowNum(ctx context.Context, uid uint64) (*UserFollowNum, error) {
	var userFollowNum UserFollowNum // 使用单个结构体实例来接收数据
	sql := fmt.Sprintf("SELECT u.id AS user_id, (SELECT COUNT(*) FROM user_follow WHERE user_id = u.id AND status = 1) AS follow_count, (SELECT COUNT(*) FROM user_follow WHERE followed_user_id = u.id AND status = 1) AS follower_count FROM user u WHERE u.id = ?;")
	err := m.QueryRowNoCacheCtx(ctx, &userFollowNum, sql, uid) // 使用 QueryRowCtx 而不是 QueryRowsNoCacheCtx
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return &userFollowNum, nil
}
