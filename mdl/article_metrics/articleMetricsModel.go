package article_metrics

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleMetricsModel = (*customArticleMetricsModel)(nil)

type (
	// ArticleMetricsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleMetricsModel.
	ArticleMetricsModel interface {
		articleMetricsModel
		FindOneWithArticleId(ctx context.Context, articleId int64) (*ArticleMetrics, error)
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

func (m *customArticleMetricsModel) FindOneWithArticleId(ctx context.Context, articleId int64) (*ArticleMetrics, error) {
	publicArticleMetricsArticleIdKey := fmt.Sprintf("cache:keyframe:am:%v", articleId)
	var resp ArticleMetrics
	err := m.QueryRowIndexCtx(ctx, &resp, publicArticleMetricsArticleIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where article_id = $1 limit 1", articleMetricsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, articleId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
