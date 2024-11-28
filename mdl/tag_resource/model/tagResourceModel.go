package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TagResourceModel = (*customTagResourceModel)(nil)

type (
	// TagResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTagResourceModel.
	TagResourceModel interface {
		tagResourceModel
	}

	customTagResourceModel struct {
		*defaultTagResourceModel
	}
)

// NewTagResourceModel returns a model for the database table.
func NewTagResourceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TagResourceModel {
	return &customTagResourceModel{
		defaultTagResourceModel: newTagResourceModel(conn, c, opts...),
	}
}
