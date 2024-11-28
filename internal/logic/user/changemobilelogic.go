package user

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更改手机号码
func NewChangeMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeMobileLogic {
	return &ChangeMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeMobileLogic) ChangeMobile(req *types.ChangeMobileRequest) (resp *types.ChangeMobileResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
