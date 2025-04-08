package user_action

import (
	"context"
	"database/sql"
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
		// ToggleAction 点赞/收藏/关注/取消点赞/取消收藏/取消关注
		// actionType 1-动态点赞, 2-评论点赞, 3-收藏 4-关注列表
		// uid 用户ID
		// targetId 目标ID
		// targetType 目标类型 1-文章 2-评论 3-用户
		ToggleAction(ctx context.Context, uid, targetId, actionType, targetType int64) error
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
	msql := fmt.Sprintf(`
    SELECT 
        article.id, 
        article.title, 
        article.author_id, 
        "user".nickname, 
        "user".avatar, 
        media.cover_url, 
        media.height, 
        media.width, 
        user_action.update_time as publish_time,

        COALESCE(article_metrics.likes, 0) AS like_count
    FROM 
        user_action
    LEFT JOIN 
        article ON user_action.target_id = article.id
    LEFT JOIN 
        "user" ON article.author_id = "user".id
    LEFT JOIN 
        media ON article.id = media.article_id
    LEFT JOIN 
        article_metrics ON article.id = article_metrics.article_id
    WHERE 
        user_action.user_id = $1 
        AND article.status = 0 
        AND user_action.target_type = 1 
        AND user_action.action_type = 1 
        AND user_action.action_value = 1
    ORDER BY 
        article.publish_time DESC 
    LIMIT 10 OFFSET $2;
`)
	err := m.QueryRowsNoCacheCtx(ctx, &list, msql, uid, offset)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

func (m *defaultUserActionModel) GetUserCollectList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error) {
	var list []*article.Feeds
	msql := fmt.Sprintf(`
    SELECT 
        article.id, 
        article.title, 
        article.author_id,
        "user".nickname, 
        "user".avatar, 
        media.cover_url, 
        media.height, 
        media.width,
        user_action.update_time as publish_time,
        COALESCE(article_metrics.likes, 0) AS like_count
    FROM 
        user_action
    LEFT JOIN 
        article ON user_action.target_id = article.id
    LEFT JOIN 
        "user" ON article.author_id = "user".id
    LEFT JOIN 
        media ON article.id = media.article_id
    LEFT JOIN 
        article_metrics ON article.id = article_metrics.article_id
    WHERE 
        user_action.user_id = $1 
        AND article.status = 0 
        AND user_action.action_type = 3 
        AND user_action.action_value = 1
    ORDER BY 
        article.publish_time DESC 
    LIMIT 10 OFFSET $2;`)

	err := m.QueryRowsNoCacheCtx(ctx, &list, msql, uid, offset)
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
	msql := `SELECT
        COALESCE(
        (
            SELECT STRING_AGG(target_id::text, ',')
            FROM (
                SELECT target_id
                FROM user_action
                WHERE user_id = $1
                  AND action_type = 1
                  AND action_value = 1
                ORDER BY update_time DESC
            ) subquery_1
        ),
        '0'
    ) AS type_1_list,
    COALESCE(
        (
            SELECT STRING_AGG(target_id::text, ',')
            FROM (
                SELECT target_id
                FROM user_action
                WHERE user_id = $1
                  AND action_type = 2
                  AND action_value = 1
                ORDER BY update_time DESC
            ) subquery_2
        ),
        '0'
    ) AS type_2_list,
    COALESCE(
        (
            SELECT STRING_AGG(target_id::text, ',')
            FROM (
                SELECT target_id
                FROM user_action
                WHERE user_id = $1
                  AND action_type = 3
                  AND action_value = 1
                ORDER BY update_time DESC
            ) subquery_3
        ),
        '0'
    ) AS type_3_list,
    COALESCE(
        (
            SELECT STRING_AGG(target_id::text, ',')
            FROM (
                SELECT target_id
                FROM user_action
                WHERE user_id = $1
                  AND action_type = 4
                  AND action_value = 1
                ORDER BY update_time DESC
            ) subquery_4
        ),
        '0'
    ) AS type_4_list;`
	keyUserActions := fmt.Sprintf("cache:keyframe:user:id:%d:action", uid)
	err := m.QueryRowCtx(ctx, &list, keyUserActions, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, v, msql, uid)
	})
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return &list, nil
}

func (m *defaultUserActionModel) ToggleAction(ctx context.Context, uid, targetId, actionType, targetType int64) error {
	keyUserActions := fmt.Sprintf("cache:keyframe:user:id:%d:action", uid)

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		msql := fmt.Sprintf(`
		INSERT INTO user_action (user_id, target_id, action_type, target_type, action_value)
		VALUES ($1, $2, $3, $4, 1)
		ON CONFLICT (user_id, target_id, action_type, target_type)
		DO UPDATE
		SET action_value = CASE WHEN user_action.action_value = 1 THEN 0 ELSE 1 END,
    		update_time = CURRENT_TIMESTAMP
		WHERE user_action.user_id = $1 AND user_action.target_id = $2 AND user_action.action_type = $3 AND user_action.target_type = $4;    
		`)
		return conn.ExecCtx(ctx, msql, uid, targetId, actionType, targetType)
	}, keyUserActions)
	if err != nil {
		logx.Error("操作失败", err)
		return err
	}
	return nil
}
