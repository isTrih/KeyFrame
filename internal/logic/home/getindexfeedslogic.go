package home

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIndexFeedsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取首页信息流
func NewGetIndexFeedsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIndexFeedsLogic {
	return &GetIndexFeedsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIndexFeedsLogic) GetIndexFeeds(req *types.GetIndexFeedsRequest) (resp *types.GetIndexFeedsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
