package media

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MediaModel = (*customMediaModel)(nil)

type (
	// MediaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMediaModel.
	MediaModel interface {
		mediaModel
		FindOneByArticleId(ctx context.Context, id uint64) (*Media, error)
	}

	customMediaModel struct {
		*defaultMediaModel
	}
)

// NewMediaModel returns a model for the database table.
func NewMediaModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MediaModel {
	return &customMediaModel{
		defaultMediaModel: newMediaModel(conn, c, opts...),
	}
}

func (m *defaultMediaModel) FindOneByArticleId(ctx context.Context, id uint64) (*Media, error) {
	keyframeMediaIdKey := fmt.Sprintf("%s%v", "cache:keyframe:media:article_id:", id)
	var resp Media
	err := m.QueryRowCtx(ctx, &resp, keyframeMediaIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `article_id` = ? limit 1", mediaRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
