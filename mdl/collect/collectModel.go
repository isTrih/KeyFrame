package collect

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CollectModel = (*customCollectModel)(nil)

type (
	// CollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCollectModel.
	CollectModel interface {
		collectModel
		GetUserCollectList(ctx context.Context, uid uint64) ([]uint64, error)
	}

	customCollectModel struct {
		*defaultCollectModel
	}
)

// NewCollectModel returns a model for the database table.
func NewCollectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CollectModel {
	return &customCollectModel{
		defaultCollectModel: newCollectModel(conn, c, opts...),
	}
}

// TODO 这里需要优化一下
func (m *defaultCollectModel) GetUserCollectList(ctx context.Context, uid uint64) ([]uint64, error) {
	chaozjCollectIdKey := fmt.Sprintf("%s%v", cacheChaozjCollectIdPrefix, uid)
	var resp []uint64
	row := "collected_feed_id"
	err := m.QueryRowCtx(ctx, &resp, chaozjCollectIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and follow_status = 1", row, m.table)
		return conn.QueryRowCtx(ctx, v, query, uid)
	})
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
