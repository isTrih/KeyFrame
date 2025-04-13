package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchWikiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSearchWikiLogic // 搜索 Wiki
func NewSearchWikiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchWikiLogic {
	return &SearchWikiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchWikiLogic) SearchWiki(req *types.SearchWikiRequest) (resp *types.WikiListResponse, err error) {
	// 设置默认值
	limit := int64(20)
	offset := int64(0)

	if req.Limit > 0 {
		limit = int64(req.Limit)
	}
	if req.Offset > 0 {
		offset = int64(req.Offset)
	}

	// 使用WikiModel搜索Wiki（现已支持全文检索）
	wikiItems, total, err := l.svcCtx.WikiModel.SearchWiki(l.ctx, req.Keyword, limit, offset)
	if err != nil {
		l.Logger.Errorf("搜索Wiki失败: %v", err)
		return nil, err
	}

	resp = &types.WikiListResponse{
		Status: "success",
		Wikis:  make([]types.WikiListItem, 0, len(wikiItems)),
		Total:  uint64(total),
	}

	// 将查询结果转换为API响应格式
	for _, item := range wikiItems {
		resp.Wikis = append(resp.Wikis, types.WikiListItem{
			Id:    uint64(item.Id),
			Title: item.Title,
		})
	}

	return resp, nil
}
