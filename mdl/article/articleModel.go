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
		GetUserFeeds(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error)
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
		Views    uint32 `db:"view_num"`
		Likes    uint32 `db:"like_num"`
	}

	FeedDetail struct {
		Id         uint64 `db:"id"`          // ID
		Title      string `db:"title"`       // 标题
		Content    string `db:"content"`     // 内容
		AuthorId   uint64 `db:"author_id"`   // 作者ID User表
		UserName   string `db:"nickname"`    // 作者昵称 User表
		Avatar     string `db:"avatar"`      // 作者头像 User表
		LikeNum    uint64 `db:"like_num"`    // 点赞数
		ShareNum   uint64 `db:"share_num"`   // 分享数
		ViewNum    uint64 `db:"view_num"`    // 浏览数
		CollectNum uint64 `db:"collect_num"` // 收藏数
		CommentNum uint64 `db:"comment_num"` // 评论数

		CoverUrl  string         `db:"cover_url"`  // 封面 Media表
		Height    uint32         `db:"height"`     // 高度 Media表
		Width     uint32         `db:"width"`      // 宽度 Media表
		MediaList sql.NullString `db:"media_list"` // 多媒体列表 Media表

		Insp       uint64 `db:"insp"`         // 人工校验为0时可用，默认为1待检验
		AiInsp     uint64 `db:"ai_insp"`      // AI检验，默认0没有问题，1出现问题。
		AiInspCode uint64 `db:"ai_insp_code"` // AI检验的状态码

		Type        uint8     `db:"type"`         // 状态 0:默认图片帧 1:视频帧 2:纯文字帧
		PublishTime time.Time `db:"publish_time"` // 发布时间
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 最后修改时间

	}
	FeedDetail2 struct {
		Id          uint64   `db:"id"`
		Title       string   `db:"title"`
		AuthorId    uint64   `db:"author_id"`
		Nickname    string   `db:"nickname"`
		Avatar      string   `db:"avatar"`
		CoverUrl    string   `db:"cover_url"`
		Height      uint32   `db:"height"`
		Width       uint32   `db:"width"`
		MediaList   []string `db:"media_list"`
		Content     string   `db:"content"`
		CreateTime  uint64   `db:"create_time"`
		UpdateTime  uint64   `db:"update_time"`
		PublishTime uint64   `db:"publish_time"`
		CollectNum  uint64   `db:"collect_num"`
		CommentNum  uint64   `db:"comment_num"`
		ShareNum    uint64   `db:"share_num"`
		AiInsp      uint64   `db:"ai_insp"`
		AiInspCode  uint64   `db:"ai_insp_code"`
		Insp        uint64   `db:"insp"`
		LikeNum     uint64   `db:"like_num"`
		ViewNum     uint64   `db:"view_num"`
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
	sql := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.status = 0 order by a.update_time desc limit 10 offset %d", row, m.table, offset)
	if len(query) != 0 {
		sql = fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.status = 0 and %s order by a.update_time desc limit 10 offset %d", row, m.table, query, offset)
	}

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

func (m *defaultArticleModel) FindOneMix(ctx context.Context, id uint64) (*Feeds, error) {
	chaozjArticleIdKey := fmt.Sprintf("%s%v", "cache:keyframe:article:preview:id:", id)
	row := "a.id,a.title,a.author_id,a.like_num,a.view_num,b.nickname,b.avatar,c.cover_url,c.height,c.width"

	var resp Feeds
	err := m.QueryRowCtx(ctx, &resp, chaozjArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.id = ? and a.status = 0 limit 1", row, m.table)
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

func (m *defaultArticleModel) FindOneDetail(ctx context.Context, id uint64) (*FeedDetail, error) {
	chaozjArticleIdKey := fmt.Sprintf("%s%v", "cache:keyframe:article:detail:id:", id)
	row := "a.id,a.title,a.content,a.author_id,b.avatar,b.nickname,a.like_num,a.share_num,a.view_num,a.collect_num,a.comment_num,c.cover_url,c.height,c.width,c.media_list,a.insp,a.ai_insp,a.ai_insp_code,a.type,a.update_time,a.publish_time,a.create_time"
	var resp FeedDetail
	err := m.QueryRowCtx(ctx, &resp, chaozjArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s as a left join user as b on a.author_id = b.id left join media as c on a.id=c.article_id where a.id = ? and a.status = 0 limit 1", row, m.table)
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
