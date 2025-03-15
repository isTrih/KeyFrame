package user_action

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zerobackend/mdl/article"
)

var _ UserActionModel = (*customUserActionModel)(nil)

type (
	// UserActionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserActionModel.
	UserActionModel interface {
		userActionModel
		GetUserCollectList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error)
		GetUserLikeList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error)
		// GetUserActionHistory 获取用户行为包括点赞收藏等行为列表
		// actionType 1-动态点赞, 2-评论点赞, 3-收藏 4-关注列表
		// uid 用户ID
		GetUserActionHistory(ctx context.Context, uid uint64) (*ActionList, error)
	}

	customUserActionModel struct {
		*defaultUserActionModel
	}

	ActionList struct {
		LikeFeedList    string `db:"type_1_list"` // 点赞帖子
		CollectFeedList string `db:"type_3_list"` // 收藏帖子
		LikeCommentList string `db:"type_2_list"` // 点赞评论
		FollowList      string `db:"type_4_list"` // 关注列表
	}
)

// NewUserActionModel returns a model for the database table.
func NewUserActionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserActionModel {
	return &customUserActionModel{
		defaultUserActionModel: newUserActionModel(conn, c, opts...),
	}
}

func (m *defaultUserActionModel) GetUserLikeList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error) {
	var list []*article.Feeds
	sql := fmt.Sprintf(`
    SELECT 
        article.id, 
        article.title, 
        article.author_id, 
        user.nickname, 
        user.avatar, 
        media.cover_url, 
        media.height, 
        media.width, 
        user_action.update_time as publish_time,

        COALESCE(action_count.like_count, 0) AS like_count
    FROM 
        user_action
    LEFT JOIN 
        article ON user_action.target_id = article.id
    LEFT JOIN 
        user ON article.author_id = user.id
    LEFT JOIN 
        media ON article.id = media.article_id
    LEFT JOIN 
        action_count ON article.id = action_count.target_id AND action_count.target_type = 1
    WHERE 
        user_action.user_id = ? 
        AND article.status = 0 
        AND user_action.target_type = 1 
        AND user_action.action_type = 1 
        AND user_action.action_value = 1
    ORDER BY 
        article.publish_time DESC 
    LIMIT 10 OFFSET ?;
`)
	err := m.QueryRowsNoCacheCtx(ctx, &list, sql, uid, offset)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

func (m *defaultUserActionModel) GetUserCollectList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error) {
	var list []*article.Feeds
	sql := fmt.Sprintf(`
    SELECT 
        article.id, 
        article.title, 
        article.author_id,
        user.nickname, 
        user.avatar, 
        media.cover_url, 
        media.height, 
        media.width,
        user_action.update_time as publish_time,
        COALESCE(action_count.like_count, 0) AS like_count
    FROM 
        user_action
    LEFT JOIN 
        article ON user_action.target_id = article.id
    LEFT JOIN 
        user ON article.author_id = user.id
    LEFT JOIN 
        media ON article.id = media.article_id
    LEFT JOIN 
        action_count ON article.id = action_count.target_id AND action_count.target_type = 1
    WHERE 
        user_action.user_id = ? 
        AND article.status = 0 
        AND user_action.action_type = 3 
        AND user_action.action_value = 1
    ORDER BY 
        article.publish_time DESC 
    LIMIT 10 OFFSET ?;
`)

	err := m.QueryRowsNoCacheCtx(ctx, &list, sql, uid, offset)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

// GetUserActionHistory 获取用户行为包括点赞收藏等行为列表
// actionType 1-动态点赞, 2-评论点赞, 3-收藏 4-关注列表
// uid 用户ID
func (m *defaultUserActionModel) GetUserActionHistory(ctx context.Context, uid uint64) (*ActionList, error) {
	var list ActionList
	sql := fmt.Sprintf(`
SELECT
    IFNULL(
        (
            SELECT GROUP_CONCAT(target_id)
            FROM user_action
            WHERE user_id = ?
              AND action_type = 1
              AND action_value = 1
            ORDER BY update_time DESC
        ),
        '0'
    ) AS type_1_list,
    IFNULL(
        (
            SELECT GROUP_CONCAT(target_id)
            FROM user_action
            WHERE user_id = ?
              AND action_type = 2
              AND action_value = 1
            ORDER BY update_time DESC
        ),
        '0'
    ) AS type_2_list,
    IFNULL(
        (
            SELECT GROUP_CONCAT(target_id)
            FROM user_action
            WHERE user_id = ?
              AND action_type = 3
              AND action_value = 1
            ORDER BY update_time DESC
        ),
        '0'
    ) AS type_3_list,
    IFNULL(
        (
            SELECT GROUP_CONCAT(target_id)
            FROM user_action
            WHERE user_id = ?
              AND action_type = 4
              AND action_value = 1
            ORDER BY update_time DESC
        ),
        '0'
    ) AS type_4_list;
`)
	keyUserActions := fmt.Sprintf("%s%v", "user:action:id:", uid)
	err := m.QueryRowCtx(ctx, &list, keyUserActions, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, v, sql, uid, uid, uid, uid)
	})
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return &list, nil
}
