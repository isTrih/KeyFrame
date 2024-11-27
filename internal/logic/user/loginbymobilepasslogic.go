package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/x/errors"
	"time"
	"zerobackend/internal/service"
	"zerobackend/mdl/user/model"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobilePassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginByMobilePassLogic 手机密码登录
func NewLoginByMobilePassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobilePassLogic {
	return &LoginByMobilePassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByMobilePassLogic) LoginByMobilePass(req *types.LoginMobilePassRequest) (resp *types.LoginResponse, err error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")

	}
	payloads := make(map[string]any)
	payloads["UID"] = user.Id
	payloads["UTYPE"] = user.Type

	accessToken, tokenErr := service.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}

	if user.Password != EncryptPassword(req.Password) {
		return nil, errors.New(6032, "密码错误")
	}

	//返回正确的token
	resp = new(types.LoginResponse)
	resp.Token = accessToken
	resp.UserId = int64(user.Id)
	return resp, nil
}
