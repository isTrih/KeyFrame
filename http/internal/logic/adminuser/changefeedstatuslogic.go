package adminuser

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeFeedStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangeFeedStatusLogic // 改变用户权限
func NewChangeFeedStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeFeedStatusLogic {
	return &ChangeFeedStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeFeedStatusLogic) ChangeFeedStatus(req *types.ChangeUserTypeRequest) (resp *types.AdminRes, err error) {
	// todo: add your logic here and delete this line

	return
}
