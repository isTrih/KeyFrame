package adminfeed

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewBanFeedLogic // 封禁动态
func NewBanFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanFeedLogic {
	return &BanFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BanFeedLogic) BanFeed(req *types.ChangeUserTypeRequest) (resp *types.AdminRes, err error) {
	// todo: add your logic here and delete this line

	return
}
