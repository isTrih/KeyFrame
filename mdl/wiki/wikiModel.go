package wiki

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WikiModel = (*customWikiModel)(nil)

type (
	// WikiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWikiModel.
	WikiModel interface {
		wikiModel
		SearchWiki(ctx context.Context, keyword string, limit, offset int64) ([]WikiListItem, int64, error)
		GetLatestWikis(ctx context.Context, limit int64) ([]WikiListItem, int64, error)
		InsertWikiHistory(ctx context.Context, wikiId, editorId int64, rawContent, changeLog string) error
	}

	customWikiModel struct {
		*defaultWikiModel
	}

	// WikiListItem 用于返回Wiki列表项
	WikiListItem struct {
		Id    int64  `db:"id"`
		Name  string `db:"name"`
		Title string // 解析后的中文标题
	}
)

// NewWikiModel returns a model for the database table.
func NewWikiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WikiModel {
	return &customWikiModel{
		defaultWikiModel: newWikiModel(conn, c, opts...),
	}
}

// SearchWiki 搜索Wiki并返回结果列表和总数
func (m *customWikiModel) SearchWiki(ctx context.Context, keyword string, limit, offset int64) ([]WikiListItem, int64, error) {
	var total int64
	baseWhereClause := "status = 1"

	if keyword == "" {
		// 无关键词时，只返回最新的Wiki列表
		return m.getWikiList(ctx, baseWhereClause, limit, offset)
	}

	// 多语言搜索条件
	multiLangSearch := `
		to_tsvector('zh_cn', raw_content) @@ to_tsquery('zh_cn', $1) OR 
		name::text ILIKE '%"zh":"%' || $1 || '%"%' OR
		name::text ILIKE '%"kr":"%' || $1 || '%"%' OR
		name::text ILIKE '%"jp":"%' || $1 || '%"%' OR
		name::text ILIKE '%"en":"%' || $1 || '%"%'
	`

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM "public"."wiki" WHERE ` + baseWhereClause + ` 
		AND (` + multiLangSearch + `)`
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 查询结果，添加排序依据
	listQuery := `
		SELECT id, name, 
		ts_rank(to_tsvector('zh_cn', raw_content), to_tsquery('zh_cn', $1)) AS rank
		FROM "public"."wiki" 
		WHERE ` + baseWhereClause + ` 
		AND (` + multiLangSearch + `)
		ORDER BY rank DESC, update_time DESC 
		LIMIT $2 OFFSET $3
	`

	var items []struct {
		Id   int64   `db:"id"`
		Name string  `db:"name"`
		Rank float64 `db:"rank"`
	}

	err = m.QueryRowsNoCacheCtx(ctx, &items, listQuery, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// 转换结果
	result := make([]WikiListItem, len(items))
	for i, item := range items {
		result[i] = WikiListItem{
			Id:   item.Id,
			Name: item.Name,
		}

		// 解析name字段的JSON
		var nameMap map[string]string
		if err := json.Unmarshal([]byte(item.Name), &nameMap); err == nil {
			result[i].Title = nameMap["zh"]
		}
	}

	return result, total, nil
}

// getWikiList 获取Wiki列表的通用方法
func (m *customWikiModel) getWikiList(ctx context.Context, whereClause string, limit, offset int64) ([]WikiListItem, int64, error) {
	// 查询总数
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "public"."wiki" WHERE %s`, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// 查询结果
	listQuery := fmt.Sprintf(`
		SELECT id, name FROM "public"."wiki" 
		WHERE %s
		ORDER BY update_time DESC 
		LIMIT %d OFFSET %d
	`, whereClause, limit, offset)

	var items []WikiListItem
	err = m.QueryRowsNoCacheCtx(ctx, &items, listQuery)
	if err != nil {
		return nil, 0, err
	}

	// 处理结果，解析name字段的JSON
	for i := range items {
		var nameMap map[string]string
		if err := json.Unmarshal([]byte(items[i].Name), &nameMap); err == nil {
			items[i].Title = nameMap["zh"]
		}
	}

	return items, total, nil
}

// GetLatestWikis 获取最新Wiki列表
func (m *customWikiModel) GetLatestWikis(ctx context.Context, limit int64) ([]WikiListItem, int64, error) {
	// 查询总数
	var total int64
	countQuery := `SELECT COUNT(*) FROM "public"."wiki" WHERE status = 1`
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// 查询最新列表
	query := fmt.Sprintf(`
		SELECT id, name FROM "public"."wiki" 
		WHERE status = 1
		ORDER BY update_time DESC 
		LIMIT %d
	`, limit)

	var items []WikiListItem
	err = m.QueryRowsNoCacheCtx(ctx, &items, query)
	if err != nil {
		return nil, 0, err
	}

	// 处理结果，解析name字段的JSON
	for i := range items {
		var nameMap map[string]string
		if err := json.Unmarshal([]byte(items[i].Name), &nameMap); err == nil {
			items[i].Title = nameMap["zh"]
		}
	}

	return items, total, nil
}

// InsertWikiHistory 插入Wiki历史记录
func (m *customWikiModel) InsertWikiHistory(ctx context.Context, wikiId, editorId int64, rawContent, changeLog string) error {
	query := `
		INSERT INTO "public"."wiki_history" (wiki_id, raw_content, editor_id, change_log) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := m.ExecNoCacheCtx(ctx, query, wikiId, rawContent, editorId, changeLog)
	return err
}
