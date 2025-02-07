package article

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"time"

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
		GetUserUploadList(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error)
		GetFeedsNum(ctx context.Context, uid uint64) (int, error)
		FindOneMix(ctx context.Context, id uint64) (*Feeds, error)
		FindOneDetail(ctx context.Context, id uint64) (*FeedDetail, error)
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
		LikeNum  uint32 `db:"like_count"`
	}

	FeedDetail struct {
		Id       uint64 `db:"id"`        // ID
		Title    string `db:"title"`     // 标题
		Content  string `db:"content"`   // 内容
		AuthorId uint64 `db:"author_id"` // 作者ID User表
		UserName string `db:"nickname"`  // 作者昵称 User表
		Avatar   string `db:"avatar"`    // 作者头像 User表

		CoverUrl  string         `db:"cover_url"`  // 封面 Media表
		Height    uint32         `db:"height"`     // 高度 Media表
		Width     uint32         `db:"width"`      // 宽度 Media表
		MediaList sql.NullString `db:"media_list"` // 多媒体列表 Media表

		PublishTime time.Time `db:"publish_time"` // 发布时间
		IpLocation  string    `db:"ip_location"`  // IP归属地
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
	row := "a.id,a.title,a.author_id,b.nickname,b.avatar,c.cover_url,c.height,c.width, COALESCE(action_count.like_count, 0) AS like_count"
	q := fmt.Sprintf(`
    SELECT %s
    FROM %s AS a
    LEFT JOIN user AS b ON a.author_id = b.id
    LEFT JOIN media AS c ON a.id = c.article_id
    LEFT JOIN action_count ON a.id = action_count.target_id AND action_count.target_type = 1
    WHERE a.status = 0
    ORDER BY a.update_time DESC
    LIMIT 10 OFFSET ?;
`, row, m.table)
	if len(query) != 0 {
		q = fmt.Sprintf(`
    SELECT %s    
    FROM %s AS a
    LEFT JOIN user AS b ON a.author_id = b.id
    LEFT JOIN media AS c ON a.id = c.article_id
    LEFT JOIN action_count ON a.id = action_count.target_id AND action_count.target_type = 1
    WHERE a.status = 0 AND (a.title LIKE '%%%s%%' OR a.content LIKE '%%%s%%' OR b.nickname LIKE '%%%s%%')
    ORDER BY a.update_time DESC
    LIMIT 10 OFFSET %d;
`, row, m.table, query, query, query, offset)

		err := m.QueryRowsNoCacheCtx(ctx, &list, q)
		if err != nil {
			logx.Error("query failed query ", err)
			return nil, err
		}
		return list, nil
	}

	err := m.QueryRowsNoCacheCtx(ctx, &list, q, offset)
	if err != nil {
		logx.Error("query failed ", err)
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

// GetFeedsNum 获取用户文章数
func (m *defaultArticleModel) GetFeedsNum(ctx context.Context, uid uint64) (int, error) {

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE `author_id` = %d", m.table, uid)

	// 定义变量来存储查询结果
	var articleCount int

	// 执行查询
	err := m.QueryRowNoCacheCtx(ctx, &articleCount, query)
	if err != nil {
		logx.Error("query failed", err)
		return 0, err
	}
	return articleCount, nil

}

// GetUserUploadList 获取用户主页文章列表
func (m *defaultArticleModel) GetUserUploadList(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error) {
	var list []*Feeds
	row := "a.id, a.title, a.author_id, b.nickname, b.avatar, c.cover_url, c.height, c.width, COALESCE(d.like_count, 0) AS like_count"
	s := fmt.Sprintf(`
    SELECT %s
    FROM %s AS a
    LEFT JOIN user AS b ON a.author_id = b.id
    LEFT JOIN media AS c ON a.id = c.article_id
    LEFT JOIN action_count AS d ON a.id = d.target_id AND d.target_type = 1
    WHERE b.id = %d
    ORDER BY a.update_time DESC
    LIMIT 10 OFFSET %d
`, row, m.table, uid, offset)

	err := m.QueryRowsNoCacheCtx(ctx, &list, s)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

// FindOneMix 根据ID获取文章预览信息
func (m *defaultArticleModel) FindOneMix(ctx context.Context, id uint64) (*Feeds, error) {
	chaozjArticleIdKey := fmt.Sprintf("%s%v", "cache:keyframe:article:preview:id:", id)
	row := "a.id, a.title, a.author_id, b.nickname, b.avatar, c.cover_url, c.height, c.width, COALESCE(d.like_count, 0) AS like_count"

	var resp Feeds
	err := m.QueryRowCtx(ctx, &resp, chaozjArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id LEFT JOIN action_count AS d ON a.id = d.target_id AND d.target_type = 1 where a.id = ? and a.status = 0 limit 1", row, m.table)
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

// FindOneDetail 获取文章详情
func (m *defaultArticleModel) FindOneDetail(ctx context.Context, id uint64) (*FeedDetail, error) {
	chaozjArticleIdKey := fmt.Sprintf("%s%v", "cache:keyframe:article:detail:id:", id)
	qury := fmt.Sprintf("SELECT article.id,article.title,article.content,article.author_id,article.ip_location,article.publish_time,user.avatar,user.nickname,media.cover_url,media.height,media.width,media.media_list FROM article LEFT JOIN user ON article.author_id = user.id LEFT JOIN media ON article.id = media.article_id WHERE article_id = ? AND article.status = 0 LIMIT 1")

	var resp FeedDetail
	err := m.QueryRowCtx(ctx, &resp, chaozjArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, v, qury, id)
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
