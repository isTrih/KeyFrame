package comment_metrics

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentMetricsModel = (*customCommentMetricsModel)(nil)

type (
	// CommentMetricsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentMetricsModel.
	CommentMetricsModel interface {
		commentMetricsModel
	}

	customCommentMetricsModel struct {
		*defaultCommentMetricsModel
	}
)

// NewCommentMetricsModel returns a model for the database table.
func NewCommentMetricsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentMetricsModel {
	return &customCommentMetricsModel{
		defaultCommentMetricsModel: newCommentMetricsModel(conn, c, opts...),
	}
}
