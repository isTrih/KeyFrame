package collect

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zerobackend/mdl/article"
)

var _ CollectModel = (*customCollectModel)(nil)

type (
	// CollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCollectModel.
	CollectModel interface {
		collectModel
		GetUserCollectList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error)
	}

	customCollectModel struct {
		*defaultCollectModel
	}
	CollectId struct {
		CollectedFeedId uint64 `db:"collected_feed_id"` // 被收藏帖子ID
	}
)

// NewCollectModel returns a model for the database table.
func NewCollectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CollectModel {
	return &customCollectModel{
		defaultCollectModel: newCollectModel(conn, c, opts...),
	}
}

func (m *defaultCollectModel) GetUserCollectList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error) {
	var list []*article.Feeds
	sql := fmt.Sprintf("SELECT article.id,article.title,article.author_id,article.like_num,article.view_num,user.nickname,user.avatar,media.cover_url,media.height,media.width FROM collect  LEFT JOIN article ON collect.collected_feed_id = article.id  LEFT JOIN user ON article.author_id = user.id  LEFT JOIN media ON article.id=media.article_id  WHERE collect.user_id = ? AND collect.follow_status = 1 AND article.status = 0 ORDER BY article.update_time DESC  LIMIT 10  OFFSET ?")
	err := m.QueryRowsNoCacheCtx(ctx, &list, sql, uid, offset)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}
