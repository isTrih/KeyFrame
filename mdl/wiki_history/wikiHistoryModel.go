package wiki_history

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WikiHistoryModel = (*customWikiHistoryModel)(nil)

type (
	// WikiHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWikiHistoryModel.
	WikiHistoryModel interface {
		wikiHistoryModel
	}

	customWikiHistoryModel struct {
		*defaultWikiHistoryModel
	}
)

// NewWikiHistoryModel returns a model for the database table.
func NewWikiHistoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WikiHistoryModel {
	return &customWikiHistoryModel{
		defaultWikiHistoryModel: newWikiHistoryModel(conn, c, opts...),
	}
}
