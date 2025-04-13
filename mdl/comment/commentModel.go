package comment

import (
	"context"
	"database/sql"
	"fmt"
	"zerobackend/internal/config"
	"zerobackend/internal/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

var (
	publicArticleMetricsArticleIdPrefix = ":public:articleMetrics:articleId:"
	publicCommentMetricsCommentIdPrefix = ":public:commentMetrics:commentId:"
)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		ReplyArticle(ctx context.Context, c config.Config, articleId, userId, insp int64, content, ipLocation string) (sql.Result, error)
		ReplyComment(ctx context.Context, c config.Config, articleId, targetId, targetUserId, userId, insp int64, content, ipLocation string) (sql.Result, error)
		DeleteComment(ctx context.Context, c config.Config, id, userId int64) error
		GetCommentList(ctx context.Context, articleId, offset uint64) ([]*CommentWithUser, uint64, error)
		LikeComment(ctx context.Context, c config.Config, commentId, userId int64) error
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

func NewCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
	}
}

func (m *defaultCommentModel) ReplyArticle(ctx context.Context, c config.Config, articleId, userId, insp int64, content, ipLocation string) (sql.Result, error) {
	// 使用指定的缓存前缀
	articleMetricsCacheKey := fmt.Sprintf("%s%v", publicArticleMetricsArticleIdPrefix, articleId)
	commentListCacheKey := fmt.Sprintf(":comment:article:%d:*", articleId) // 新增评论列表缓存键模式
	commentlistKeys, _ := utils.RedisKey(c, commentListCacheKey)
	cacheKeys := []string{articleMetricsCacheKey}
	cacheKeys = append(cacheKeys, commentlistKeys...)
	query := `
        WITH new_comment AS (
            INSERT INTO "public"."comment" (
                "article_id", "parent_id", "parent_user_id", "user_id", 
                "content", "ip_location", "ai_insp"
            ) VALUES ($1, 0, -1, $2, $3, $4, $5)
            RETURNING id, article_id
        ),
        update_article_metrics AS (
            UPDATE "public"."article_metrics"
            SET "comments" = "comments" + 1
            WHERE "article_id" = (SELECT article_id FROM new_comment)
        ),
        ensure_comment_metrics AS (
            INSERT INTO "public"."comment_metrics" ("comment_id", "likes", "comments")
            SELECT id, 0, 0 FROM new_comment
            ON CONFLICT ("comment_id") DO NOTHING
        )
        SELECT id FROM new_comment
    `

	// 使用 ExecCtx 在操作完成后清除缓存
	result, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, articleId, userId, content, ipLocation, insp)
	}, cacheKeys...) // 添加评论列表缓存键
	return result, err
}

func (m *defaultCommentModel) ReplyComment(ctx context.Context, c config.Config, articleId, targetId, targetUserId, userId, insp int64, content, ipLocation string) (sql.Result, error) {
	// 使用指定的缓存前缀
	articleMetricsCacheKey := fmt.Sprintf("%s%v", publicArticleMetricsArticleIdPrefix, articleId)
	commentMetricsCacheKey := fmt.Sprintf("%s%v", publicCommentMetricsCommentIdPrefix, targetId)
	commentListCacheKey := fmt.Sprintf(":comment:article:%d:*", articleId) // 新增评论列表缓存键模式
	commentlistKeys, _ := utils.RedisKey(c, commentListCacheKey)
	cacheKeys := []string{articleMetricsCacheKey, commentMetricsCacheKey}
	cacheKeys = append(cacheKeys, commentlistKeys...)
	query := `
        WITH new_reply AS (
            INSERT INTO "public"."comment" (
                "article_id", "parent_id", "parent_user_id", "user_id", 
                "content", "ip_location", "ai_insp"
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id, article_id, parent_id
        ),
        update_article_metrics AS (
            UPDATE "public"."article_metrics"
            SET "comments" = "comments" + 1
            WHERE "article_id" = (SELECT article_id FROM new_reply)
        ),
        ensure_parent_metrics AS (
            INSERT INTO "public"."comment_metrics" ("comment_id", "likes", "comments")
            SELECT article_id, 0, 0 FROM new_reply
            WHERE article_id != 0
            ON CONFLICT ("comment_id") DO NOTHING
        ),
        update_parent_metrics AS (
            UPDATE "public"."comment_metrics"
            SET "comments" = "comments" + 1
            WHERE "comment_id" = (SELECT article_id FROM new_reply WHERE article_id != 0)
        ),
        ensure_comment_metrics AS (
            INSERT INTO "public"."comment_metrics" ("comment_id", "likes", "comments")
            SELECT id, 0, 0 FROM new_reply
            ON CONFLICT ("comment_id") DO NOTHING
        )
        SELECT id FROM new_reply
    `

	// 使用 ExecCtx 并在操作完成后清除缓存
	result, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, articleId, targetId, targetUserId, userId, content, ipLocation, insp)
	}, cacheKeys...) // 添加评论列表缓存键

	return result, err
}

func (m *defaultCommentModel) DeleteComment(ctx context.Context, c config.Config, id, userId int64) error {
	// 首先查询评论获取 article_id，用于生成缓存键
	comment, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	// 使用指定的缓存前缀
	articleMetricsCacheKey := fmt.Sprintf("%s%v", publicArticleMetricsArticleIdPrefix, comment.ArticleId)
	commentMetricsCacheKey := fmt.Sprintf("%s%v", publicCommentMetricsCommentIdPrefix, comment.ParentId)
	commentListCacheKey := fmt.Sprintf(":comment:article:%d:*", comment.ArticleId) // 新增评论列表缓存键模式
	commentlistKeys, _ := utils.RedisKey(c, commentListCacheKey)
	cacheKeys := []string{publicCommentIdPrefix + fmt.Sprintf("%v", id), articleMetricsCacheKey, commentMetricsCacheKey}
	cacheKeys = append(cacheKeys, commentlistKeys...)
	query := `
		WITH deleted AS (
			DELETE FROM "public"."comment"
			WHERE "id" = $1 AND "user_id" = $2
			RETURNING article_id, article_id
		),
		update_article_metrics AS (
			UPDATE "public"."article_metrics"
			SET "comments" = "comments" - 1
			WHERE "article_id" = (SELECT article_id FROM deleted)
		),
		update_comment_metrics AS (
			UPDATE "public"."comment_metrics"
			SET "comments" = "comments" - 1
			WHERE "comment_id" = (SELECT article_id FROM deleted WHERE article_id != 0)
		)
		SELECT 1
	`

	// 使用 ExecCtx 并在操作完成后清除缓存
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, id, userId)
	}, cacheKeys...)
	return err
}

func (m *defaultCommentModel) GetCommentList(ctx context.Context, articleId, offset uint64) ([]*CommentWithUser, uint64, error) {
	key := fmt.Sprintf(":comment:article:%d:%d", articleId, offset)

	var comments []*CommentWithUser
	err := m.QueryRowCtx(ctx, &comments, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := `
        WITH main_comments AS (
            SELECT 
                c.id, c.user_id, u.nickname, u.avatar, 
                c.content, COALESCE(cm.likes, 0) as like_count,
                EXTRACT(EPOCH FROM c.create_time AT TIME ZONE 'Asia/Shanghai')::bigint as create_time, 
                c.parent_id, c.parent_user_id,
                '' as reply_to_nickname,
                c.ip_location
            FROM comment c
            JOIN "user" u ON c.user_id = u.id
            LEFT JOIN comment_metrics cm ON c.id = cm.comment_id
            WHERE c.article_id = $1 
              AND c.parent_id = 0 
              AND c.status = 0
              AND (CASE WHEN c.ai_insp != 0 THEN c.insp = 0
                   ELSE 1 = 1 END)
            ORDER BY c.create_time DESC
            LIMIT $2 OFFSET $3
        ),
        sub_comments AS (
            SELECT 
                c.id, c.user_id, u.nickname, u.avatar, 
                c.content, COALESCE(cm.likes, 0) as like_count,
                EXTRACT(EPOCH FROM c.create_time AT TIME ZONE 'Asia/Shanghai')::bigint as create_time,
                c.parent_id, c.parent_user_id,
                target_user.nickname as reply_to_nickname,
                c.ip_location
            FROM comment c
            JOIN "user" u ON c.user_id = u.id
            JOIN main_comments parent ON c.parent_id = parent.id
            JOIN "user" target_user ON c.parent_user_id = target_user.id
            LEFT JOIN comment_metrics cm ON c.id = cm.comment_id
            WHERE c.status = 0
              AND (CASE WHEN c.ai_insp != 0 THEN c.insp = 0
                   ELSE 1 = 1 END)
            ORDER BY c.create_time DESC
        )
        SELECT * FROM main_comments
        UNION ALL
        SELECT * FROM sub_comments`
		return conn.QueryRowsCtx(ctx, v, query, articleId, 20, offset)
	})

	if len(comments) == 0 {
		err := m.DelCacheCtx(ctx, key)
		if err != nil {
			return nil, 0, err
		}
	}
	if err != nil {
		return nil, 0, err
	}

	total, err := m.getCommentTotal(ctx, articleId)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (m *defaultCommentModel) getCommentTotal(ctx context.Context, articleId uint64) (uint64, error) {
	var total uint64
	err := m.QueryRowNoCacheCtx(ctx, &total, `
        SELECT COUNT(*) FROM comment 
        WHERE article_id = $1 AND status = 0
        AND status = 0
          AND (CASE WHEN ai_insp != 0 THEN insp = 0
        ELSE 1 = 1 END)`, articleId)
	return total, err
}

func (m *defaultCommentModel) LikeComment(ctx context.Context, c config.Config, commentId, userId int64) error {
	// 首先查询评论获取 article_id，用于生成缓存键
	comment, err := m.FindOne(ctx, commentId)
	if err != nil {
		return err
	}
	commentListCacheKey := fmt.Sprintf(":comment:article:%d:*", comment.ArticleId) // 新增评论列表缓存键模式

	commentlistKeys, _ := utils.RedisKey(c, commentListCacheKey)

	commentMetricsCacheKey := fmt.Sprintf("%s%v", publicCommentMetricsCommentIdPrefix, commentId)
	userActionCacheKey := fmt.Sprintf("cache:keyframe:user:id:%d:action", userId)
	cacheKeys := []string{commentMetricsCacheKey, userActionCacheKey}
	cacheKeys = append(cacheKeys, commentlistKeys...)
	query := `
WITH action AS (
    INSERT INTO "public"."user_action" ("user_id", "target_id", "action_type", "target_type", "action_value")
    VALUES ($1, $2, 2, 2, 1)
    ON CONFLICT ("user_id", "target_id", "action_type") 
    DO UPDATE SET 
        "action_value" = CASE 
            WHEN "user_action"."action_value" = 1 THEN 0
            ELSE 1
        END
    RETURNING "action_value"
),
update_metrics AS (
    UPDATE "public"."comment_metrics"
    SET "likes" = "likes" + CASE 
        WHEN (SELECT "action_value" FROM action) = 1 THEN 1
        ELSE -1
    END
    WHERE "comment_id" = $2
)
SELECT 1;  `

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, userId, commentId)
	}, cacheKeys...)
	return err
}

type (
	CommentWithUser struct {
		Id              uint64
		UserId          uint64
		Nickname        string
		Avatar          string
		Content         string
		LikeCount       uint64
		CreateTime      uint64
		ParentId        uint64
		ParentUserId    int64
		ReplyToNickname string
		IpLocation      string
	}
)
