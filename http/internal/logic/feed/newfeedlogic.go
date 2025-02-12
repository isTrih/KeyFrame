package feed

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNewFeedLogic // 创建帧（文章）
func NewNewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewFeedLogic {
	return &NewFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewFeedLogic) NewFeed(req *types.NewFeedRequest) (resp *types.StatusResponse, err error) {
	// 流程：客户端上传图片（带有自动删除的生命周期），完成发布后更改生命周期
	// todo: add your logic here and delete this line

	return
}
