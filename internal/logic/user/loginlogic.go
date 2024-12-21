package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"time"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/internal/utils"
	model "zerobackend/mdl/user"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")

	}

	// 校验密码
	if user.Password != EncryptPassword(req.Password) {
		return nil, errors.New(6032, "密码错误")
	}

	// 生成token
	payloads := make(map[string]any)
	payloads["UID"] = user.Id
	payloads["UTYPE"] = user.Type
	payloads["USTATUS"] = user.Status
	accessToken, tokenErr := utils.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}

	// 返回token
	resp = new(types.LoginResponse)
	resp.Token = accessToken
	resp.UserType = uint8(user.Type)
	resp.Avatar = user.Avatar
	resp.Signature = user.Signature
	resp.UserName = user.Nickname
	resp.UserId = user.Id
	return resp, nil
}
