package home

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSearchFeedLogic // 搜索帧（文章）
func NewSearchFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFeedLogic {
	return &SearchFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchFeedLogic) SearchFeed(req *types.GetIndexFeedsRequest) (resp *types.GetIndexFeedsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
