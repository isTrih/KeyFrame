package wiki

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WikiModel = (*customWikiModel)(nil)

type (
	// WikiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWikiModel.
	WikiModel interface {
		wikiModel
	}

	customWikiModel struct {
		*defaultWikiModel
	}
)

// NewWikiModel returns a model for the database table.
func NewWikiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WikiModel {
	return &customWikiModel{
		defaultWikiModel: newWikiModel(conn, c, opts...),
	}
}
