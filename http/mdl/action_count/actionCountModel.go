package action_count

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ActionCountModel = (*customActionCountModel)(nil)

type (
	// ActionCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customActionCountModel.
	ActionCountModel interface {
		actionCountModel
	}

	customActionCountModel struct {
		*defaultActionCountModel
	}
)

// NewActionCountModel returns a model for the database table.
func NewActionCountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ActionCountModel {
	return &customActionCountModel{
		defaultActionCountModel: newActionCountModel(conn, c, opts...),
	}
}
