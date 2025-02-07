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
	// todo: add your logic here and delete this line

	return
}
