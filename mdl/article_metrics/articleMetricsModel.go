package article_metrics

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleMetricsModel = (*customArticleMetricsModel)(nil)

type (
	// ArticleMetricsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleMetricsModel.
	ArticleMetricsModel interface {
		articleMetricsModel
	}

	customArticleMetricsModel struct {
		*defaultArticleMetricsModel
	}
)

// NewArticleMetricsModel returns a model for the database table.
func NewArticleMetricsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleMetricsModel {
	return &customArticleMetricsModel{
		defaultArticleMetricsModel: newArticleMetricsModel(conn, c, opts...),
	}
}
