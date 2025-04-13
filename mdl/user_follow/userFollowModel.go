package user_follow

import (
	"context"
	"database/sql"
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
		ToggleFollow(ctx context.Context, uid, targetId int64) error
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
	sql := fmt.Sprintf(`
		SELECT u.id AS user_id, 
       (SELECT COUNT(*) FROM "public"."user_follow" WHERE user_id = u.id AND status = 1) AS follow_count, 
       (SELECT COUNT(*) FROM "public"."user_follow" WHERE followed_user_id = u.id AND status = 1) AS follower_count 
		FROM "public"."user" u WHERE u.id = $1`)
	err := m.QueryRowNoCacheCtx(ctx, &userFollowNum, sql, uid) // 使用 QueryRowCtx 而不是 QueryRowsNoCacheCtx
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return &userFollowNum, nil
}

func (m *defaultUserFollowModel) ToggleFollow(ctx context.Context, uid, targetId int64) error {
	keyUserActions := fmt.Sprintf("cache:keyframe:user:id:%d:action", uid)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		msql := `
        WITH update_action AS (
            UPDATE user_action 
            SET action_value = CASE action_value WHEN 1 THEN 0 ELSE 1 END,
                update_time = CURRENT_TIMESTAMP
            WHERE user_id = $1 
              AND target_id = $2 
              AND action_type = 4 
              AND target_type = 3
            RETURNING action_value
        ),
        insert_action AS (
            INSERT INTO user_action (user_id, target_id, action_type, target_type, action_value)
            SELECT $1, $2, 4, 3, 1
            WHERE NOT EXISTS (SELECT 1 FROM update_action)
            RETURNING action_value
        ),
        combined_action AS (
            SELECT action_value FROM update_action
            UNION ALL
            SELECT action_value FROM insert_action
        ),
        -- 优化点：将插入逻辑提前，确保记录存在
        upsert_follow AS (
            INSERT INTO user_follow (user_id, followed_user_id, status)
            SELECT $1, $2, (SELECT action_value FROM combined_action)
            ON CONFLICT (user_id, followed_user_id) 
            DO UPDATE SET 
                status = EXCLUDED.status,
                update_time = CURRENT_TIMESTAMP
            RETURNING status
        ),
        -- 更新操作改为基于插入结果
        update_follow AS (
            UPDATE user_follow
            SET status = (SELECT status FROM upsert_follow),
                update_time = CURRENT_TIMESTAMP
            WHERE user_id = $1 AND followed_user_id = $2
        )
        SELECT 1` // ← Add final SELECT statement

		return conn.ExecCtx(ctx, msql, uid, targetId)
	}, keyUserActions)
	if err != nil {
		logx.Error("操作失败", err)
		return err
	}
	return nil
}
