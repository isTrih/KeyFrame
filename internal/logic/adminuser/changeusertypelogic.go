package adminuser

import (
	"context"

	"zerobackend/http/internal/svc"
	"zerobackend/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeUserTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangeUserTypeLogic // 改变用户权限
func NewChangeUserTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserTypeLogic {
	return &ChangeUserTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserTypeLogic) ChangeUserType(req *types.ChangeUserTypeRequest) (resp *types.AdminRes, err error) {
	// todo: add your logic here and delete this line

	return
}
