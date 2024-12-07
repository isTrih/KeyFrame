package article

import (
	"context"
	"fmt"
	//"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		GetFeeds(ctx context.Context, offset uint64, query string) ([]*Feeds, error)
		GetUncheckFeeds(ctx context.Context, offset uint64) ([]*Feeds, error)
		GetUserFeeds(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error)
	}

	customArticleModel struct {
		*defaultArticleModel
	}
	Feeds struct {
		Id       uint64 `db:"id"`
		Title    string `db:"title"`
		AuthorId uint64 `db:"author_id"`
		UserName string `db:"nickname"`
		Avatar   string `db:"avatar"`
		CoverUrl string `db:"cover_url"`
		Height   uint32 `db:"height"`
		Width    uint32 `db:"width"`
		Views    uint32 `db:"view_num"`
		Likes    uint32 `db:"like_num"`
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}

// GetFeeds 获取文章 query是接在 Where 的 and 后面的字符串 如："a.id = 1 and a.title = 'test'"
func (m *defaultArticleModel) GetFeeds(ctx context.Context, offset uint64, query string) ([]*Feeds, error) {
	var list []*Feeds
	row := "a.id,a.title,a.author_id,a.like_num,a.view_num,b.nickname,b.avatar,c.cover_url,c.height,c.width"
	sql := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.status = 2 order by a.update_time desc limit 10 offset %d", row, m.table, offset)
	if len(query) == 0 {
		sql = fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.status = 2 and %s order by a.update_time desc limit 10 offset %d", row, m.table, query, offset)
	}

	//sqlBuilder := squirrel.Select("a.id", "a.title", "a.author_id", "b.nickname", "b.avatar")
	//sqlBuilder.From(m.table).Offset(offset).Limit(10).OrderBy("update_time DESC")
	//sqlStr, args, err := sqlBuilder.ToSql()
	//if err != nil {
	//	logx.Error("generate sqlStr failed", err)
	//	return nil, err
	//}
	err := m.QueryRowsNoCacheCtx(ctx, &list, sql)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

// GetUncheckFeeds 获取未审核文章
func (m *defaultArticleModel) GetUncheckFeeds(ctx context.Context, offset uint64) ([]*Feeds, error) {

	var list []*Feeds
	row := "a.id,a.title,a.author_id,a.like_num,a.view_num,b.nickname,b.avatar,c.cover_url,c.height,c.width"
	sql := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.status = 0 order by a.update_time desc limit 10 offset %d", row, m.table, offset)

	err := m.QueryRowsNoCacheCtx(ctx, &list, sql)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

// GetUserFeeds 获取用户主页文章
func (m *defaultArticleModel) GetUserFeeds(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error) {

	var list []*Feeds
	row := "a.id,a.title,a.author_id,a.like_num,a.view_num,b.nickname,b.avatar,c.cover_url,c.height,c.width"
	sql := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where b.id = %d order by a.update_time desc limit 10 offset %d", row, m.table, uid, offset)

	err := m.QueryRowsNoCacheCtx(ctx, &list, sql)
	if err != nil {
		return nil, err
	}
	return list, nil
}
