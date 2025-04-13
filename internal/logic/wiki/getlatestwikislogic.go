package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestWikisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLatestWikisLogic // 获取最新 Wiki 列表
func NewGetLatestWikisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestWikisLogic {
	return &GetLatestWikisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLatestWikisLogic) GetLatestWikis(req *types.GetLatestWikisRequest) (resp *types.WikiListResponse, err error) {
	// 设置查询限制
	limit := int64(req.Limit)
	if limit <= 0 {
		limit = 50 // 默认取50条
	}

	// 使用WikiModel获取最新Wiki列表
	wikiItems, total, err := l.svcCtx.WikiModel.GetLatestWikis(l.ctx, limit)
	if err != nil {
		l.Logger.Errorf("获取最新Wiki列表失败: %v", err)
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
