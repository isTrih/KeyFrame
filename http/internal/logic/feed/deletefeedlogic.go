package feed

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteFeedLogic // 删除帧（文章）
func NewDeleteFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFeedLogic {
	return &DeleteFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFeedLogic) DeleteFeed(req *types.DeleteFeedRequest) (resp *types.StatusResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
