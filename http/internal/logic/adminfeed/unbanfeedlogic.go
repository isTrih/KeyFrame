package adminfeed

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnBanFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUnBanFeedLogic // 解禁动态
func NewUnBanFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnBanFeedLogic {
	return &UnBanFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnBanFeedLogic) UnBanFeed(req *types.ChangeUserTypeRequest) (resp *types.AdminRes, err error) {
	// todo: add your logic here and delete this line

	return
}
