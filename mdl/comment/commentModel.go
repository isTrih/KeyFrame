package comment

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		ReplyArticle(ctx context.Context, articleId, userId int64, content, ipLocation string) (sql.Result, error)
		ReplyComment(ctx context.Context, articleId, targetId, targetUserId, userId int64, content, ipLocation string) (sql.Result, error)
		DeleteComment(ctx context.Context, id, userId int64) error
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
	}
}

func (m *defaultCommentModel) ReplyArticle(ctx context.Context, articleId, userId int64, content, ipLocation string) (sql.Result, error) {
	query := `
		WITH new_comment AS (
			INSERT INTO "public"."comment" (
				"article_id", "target_id", "target_user_id", "user_id", 
				"content", "ip_location"
			) VALUES ($1, 0, -1, $2, $3, $4)
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

	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, articleId, userId, content, ipLocation)
	})
}

func (m *defaultCommentModel) ReplyComment(ctx context.Context, articleId, targetId, targetUserId, userId int64, content, ipLocation string) (sql.Result, error) {
	query := `
		WITH new_reply AS (
			INSERT INTO "public"."comment" (
				"article_id", "target_id", "target_user_id", "user_id", 
				"content", "ip_location"
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, article_id, target_id
		),
		update_article_metrics AS (
			UPDATE "public"."article_metrics"
			SET "comments" = "comments" + 1
			WHERE "article_id" = (SELECT article_id FROM new_reply)
		),
		ensure_parent_metrics AS (
			INSERT INTO "public"."comment_metrics" ("comment_id", "likes", "comments")
			SELECT target_id, 0, 0 FROM new_reply
			WHERE target_id != 0
			ON CONFLICT ("comment_id") DO NOTHING
		),
		update_parent_metrics AS (
			UPDATE "public"."comment_metrics"
			SET "comments" = "comments" + 1
			WHERE "comment_id" = (SELECT target_id FROM new_reply WHERE target_id != 0)
		),
		ensure_comment_metrics AS (
			INSERT INTO "public"."comment_metrics" ("comment_id", "likes", "comments")
			SELECT id, 0, 0 FROM new_reply
			ON CONFLICT ("comment_id") DO NOTHING
		)
		SELECT id FROM new_reply
	`

	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, articleId, targetId, targetUserId, userId, content, ipLocation)
	})
}

func (m *defaultCommentModel) DeleteComment(ctx context.Context, id, userId int64) error {
	query := `
		WITH deleted AS (
			DELETE FROM "public"."comment"
			WHERE "id" = $1 AND "user_id" = $2
			RETURNING article_id, target_id
		),
		update_article_metrics AS (
			UPDATE "public"."article_metrics"
			SET "comments" = "comments" - 1
			WHERE "article_id" = (SELECT article_id FROM deleted)
		),
		update_comment_metrics AS (
			UPDATE "public"."comment_metrics"
			SET "comments" = "comments" - 1
			WHERE "comment_id" = (SELECT target_id FROM deleted WHERE target_id != 0)
		)
		SELECT 1
	`

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, id, userId)
	})
	return err
}
