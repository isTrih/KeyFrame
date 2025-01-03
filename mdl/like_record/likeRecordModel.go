package like_record

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zerobackend/mdl/article"
)

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
	// LikeRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeRecordModel.
	LikeRecordModel interface {
		likeRecordModel
		GetUserLikeList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error)
	}

	customLikeRecordModel struct {
		*defaultLikeRecordModel
	}
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeRecordModel {
	return &customLikeRecordModel{
		defaultLikeRecordModel: newLikeRecordModel(conn, c, opts...),
	}
}

func (m *defaultLikeRecordModel) GetUserLikeList(ctx context.Context, offset uint64, uid uint64) ([]*article.Feeds, error) {
	var list []*article.Feeds
	sql := fmt.Sprintf("SELECT article.id,article.title,article.author_id,article.like_num,article.view_num,user.nickname,user.avatar,media.cover_url,media.height,media.width FROM like_record LEFT JOIN article ON like_record.obj_id = article.id  LEFT JOIN user ON article.author_id = user.id  LEFT JOIN media ON article.id=media.article_id  WHERE like_record.user_id = ? AND article.status = 2 AND like_record.obj_id = 0 ORDER BY article.update_time DESC  LIMIT 10  OFFSET ?")
	err := m.QueryRowsNoCacheCtx(ctx, &list, sql, uid, offset)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}
