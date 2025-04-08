package user

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangeUserInfoLogic 更改用户信息 需要token
func NewChangeUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserInfoLogic {
	return &ChangeUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserInfoLogic) ChangeUserInfo(req *types.ChangeUserInfoRequest) (resp *types.ChangeUserInfoResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
