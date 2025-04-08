package article

import (
	"context"
	"database/sql"
	"encoding/json"
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
		GetUncheckFeeds(ctx context.Context, statusType int16, offset uint64) ([]*Feeds, error)
		GetUserUploadList(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error)
		GetFeedsNum(ctx context.Context, uid uint64) (int, error)
		FindOneMix(ctx context.Context, id uint64) (*Feeds, error)
		FindOneDetail(ctx context.Context, id uint64) (*FeedDetail, error)
		NewFeed(ctx context.Context, title string, content string, rawContent string, cover string, height uint32, width uint32, media []string, uid int64, region string, insp uint8) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
	Feeds struct {
		Id          uint64    `db:"id"`
		Title       string    `db:"title"`
		AuthorId    uint64    `db:"author_id"`
		UserName    string    `db:"nickname"`
		Avatar      string    `db:"avatar"`
		CoverUrl    string    `db:"cover_url"`
		Height      uint32    `db:"height"`
		Width       uint32    `db:"width"`
		LikeNum     uint32    `db:"like_count"`
		PublishTime time.Time `db:"publish_time"` // 发布时间
	}

	FeedDetail struct {
		Id        uint64         `db:"id"`         // ID
		Title     string         `db:"title"`      // 标题
		Content   string         `db:"content"`    // 内容
		AuthorId  uint64         `db:"author_id"`  // 作者ID User表
		UserName  string         `db:"nickname"`   // 作者昵称 User表
		Avatar    string         `db:"avatar"`     // 作者头像 User表
		Type      uint16         `db:"type"`       // 作者类型 User表
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
	row := "a.id,a.title,a.author_id,a.publish_time,b.nickname,b.avatar,c.cover_url,c.height,c.width, COALESCE(article_metrics.likes, 0) AS like_count"
	q := fmt.Sprintf(`
SELECT %s
FROM %s AS a
LEFT JOIN "user" AS b ON a.author_id = b.id
LEFT JOIN media AS c ON a.id = c.article_id
LEFT JOIN article_metrics ON a.id = article_metrics.article_id
WHERE a.status = 0   
AND 
    (CASE WHEN a.ai_insp != 0 THEN a.insp = 0
        ELSE 1 = 1  -- 当 a.ai_insp 为 0 时，此条件恒为真，即不考虑 a.insp 的值
    END)
ORDER BY a.publish_time DESC
LIMIT 18 OFFSET $1
`, row, m.table)
	if len(query) != 0 {
		// TODO:这里之后换成全文搜索
		q := `
        SELECT ` + row + `,
		ts_rank(to_tsvector('zh_cn', a.raw_content), to_tsquery('zh_cn', $1)) AS rank
        FROM ` + m.table + ` AS a
        LEFT JOIN "user" AS b ON a.author_id = b.id
        LEFT JOIN media AS c ON a.id = c.article_id
		LEFT JOIN article_metrics ON a.id = article_metrics.article_id
        WHERE 
		  (CASE WHEN a.ai_insp != 0 THEN a.insp = 0 ELSE true END)
 		  AND a.status = 0 
          AND (to_tsvector('zh_cn', raw_content) @@ to_tsquery('zh_cn', $1) OR a.title LIKE $1 || '%' OR b.nickname LIKE $1 || '%')
        ORDER BY a.publish_time DESC ,rank DESC
        LIMIT 18 OFFSET $2
    `

		err := m.QueryRowsNoCacheCtx(ctx, &list, q, query, offset)
		if err != nil {
			logx.Error("query failed query ", err)
			return nil, err
		}
		return list, nil
	}

	// 不搜索
	err := m.QueryRowsNoCacheCtx(ctx, &list, q, offset)
	if err != nil {
		logx.Error("query failed ", err)
		return nil, err
	}
	return list, nil
}

// GetUncheckFeeds 获取审核未通过的文章
func (m *defaultArticleModel) GetUncheckFeeds(ctx context.Context, statusType int16, offset uint64) ([]*Feeds, error) {

	var list []*Feeds
	row := "a.id,a.title,a.author_id, a.publish_time, a.like_num,a.view_num,b.nickname,b.avatar,c.cover_url,c.height,c.width"
	mysql := fmt.Sprintf(`
              select %s from %s as a left join "user" as b on a.author_id = b.id 
              left join "media" as c on a.id=c.article_id 
              where a.status = %s order by a.update_time desc limit 10 offset %d`,
		row, m.table, statusType, offset)

	err := m.QueryRowsNoCacheCtx(ctx, &list, mysql)
	if err != nil {
		logx.Error("query failed", err)
		return nil, err
	}
	return list, nil
}

// GetFeedsNum 获取用户文章数
func (m *defaultArticleModel) GetFeedsNum(ctx context.Context, uid uint64) (int, error) {

	query := `SELECT COUNT(*) FROM "public"."article" WHERE "author_id" = $1`

	// 定义变量来存储查询结果
	var articleCount int

	// 执行查询
	err := m.QueryRowNoCacheCtx(ctx, &articleCount, query, uid)
	if err != nil {
		logx.Error("query failed", err)
		return 0, err
	}
	return articleCount, nil

}

// GetUserUploadList 获取用户主页文章列表
func (m *defaultArticleModel) GetUserUploadList(ctx context.Context, offset uint64, uid uint64) ([]*Feeds, error) {
	var list []*Feeds
	row := "a.id, a.title, a.author_id, a.publish_time, b.nickname, b.avatar, c.cover_url, c.height, c.width, COALESCE(d.likes, 0) AS like_count"
	s := fmt.Sprintf(`
    SELECT %s
    FROM %s AS a
    LEFT JOIN "user" AS b ON a.author_id = b.id
    LEFT JOIN media AS c ON a.id = c.article_id
    LEFT JOIN article_metrics AS d ON a.id = d.article_id 
    WHERE b.id = %d        
    	AND 
        (CASE WHEN a.ai_insp != 0 THEN a.insp = 0
    		  ELSE 1 = 1  -- 当 a.ai_insp 为 0 时，此条件恒为真，即不考虑 a.insp 的值
    	END)
    ORDER BY a.publish_time DESC
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
	row := "a.id, a.title, a.author_id, a.publish_time, b.nickname, b.avatar, c.cover_url, c.height, c.width, COALESCE(d.likes, 0) AS like_count"

	var resp Feeds
	err := m.QueryRowCtx(ctx, &resp, chaozjArticleIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf(`select %s from %s as a left join "user" as b on a.author_id = b.id left join "media" as c on a.id=c.article_id LEFT JOIN article_metrics AS d ON a.id = d.article_id where a.id = $1 and a.status = 0 limit 1`,
			row, m.table)
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
	qury := fmt.Sprintf(`SELECT article.id,article.title,article.content,article.author_id,article.ip_location,article.publish_time,"user".avatar,"user".nickname,"user".type,media.cover_url,media.height,media.width,media.media_list 
FROM article 
    LEFT JOIN "user" ON article.author_id = "user".id LEFT JOIN media ON article.id = media.article_id WHERE article_id = $1 AND article.status = 0 LIMIT 1`)

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

// NewFeed 新建文章
func (m *defaultArticleModel) NewFeed(ctx context.Context,
	title string, content string, rawContent string,
	cover string, height uint32, width uint32, media []string,
	uid int64, region string, insp uint8,
) error {
	// 将切片转换为 JSON 字符串
	mediaListJSON, err := json.Marshal(media)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return err
	}

	// 文章缓存
	cachePublicArticleIdKey := fmt.Sprintf("%s%v", cachePublicArticleIdPrefix, uid)
	// 用户文章列表缓存
	cacheUserUploadListKey := fmt.Sprintf("cache:keyframe:user:id:%d:upload", uid)

	// 插入文章+媒体数据
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		mysql := `
		WITH inserted_article AS (
    	INSERT INTO article (title, content, raw_content, author_id, ip_location, publish_time, ai_insp)
    	VALUES ($1, $2, $3, $4, $5, $6, $7)
    	RETURNING id)
		INSERT INTO media (article_id, cover_url, height, width, media_list)
		SELECT id, $8, $9, $10, $11
		FROM inserted_article;`
		return conn.ExecCtx(ctx, mysql, title, content, rawContent, uid, region, time.Now(), insp, cover, height, width, mediaListJSON)
	}, cachePublicArticleIdKey, cacheUserUploadListKey)

	if err != nil {
		fmt.Println("插入表失败:", err)
		return err
	}
	return nil
}

// EditFeed 编辑文章
func (m *defaultArticleModel) EditFeed(ctx context.Context,
	title string, content string, rawContent string,
	cover string, height uint32, width uint32, media []string,
	uid int64, region string, insp uint8,
) error {
	//// 文章详情缓存
	//cacheArticleDetailKey := fmt.Sprintf("cache:keyframe:article:detail:id:%d", id)
	//// 文章预览缓存
	//cacheArticlePreviewKey := fmt.Sprintf("%s%v", "cache:keyframe:article:preview:id:", id)
	return nil
}
