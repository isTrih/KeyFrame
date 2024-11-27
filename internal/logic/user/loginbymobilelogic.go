package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/x/errors"
	"zerobackend/mdl/user/model"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户手机登录
func NewLoginByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileLogic {
	return &LoginByMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByMobileLogic) LoginByMobile(req *types.LoginByMobileRequest) (resp *types.LoginResponse, err error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")

	}

	return &types.LoginResponse{
		UserId: int64(user.Id),
		Token:  user.Nickname,
	}, nil
}
