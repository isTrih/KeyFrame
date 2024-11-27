package user

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"time"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/mdl/user/model"
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
	// todo: add your logic here and delete this line

	// 模拟登录逻辑
	if req.UserId != "go-zero" || req.Password != "123456" {
		return nil, errors.New(1001, "用户名或密码错误")
	}

	var userData model.User
	payloads := make(map[string]any)
	payloads["userIdentity"] = userData.Id

	accessToken, tokenErr := l.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return &types.LoginResponse{
		Token: accessToken,
	}, nil

	resp = new(types.LoginResponse)
	resp.Token = "go-zero"
	resp.UserId = 1
	return resp, nil
}

func (l *LoginLogic) GetToken(iat int64, secretKey string, payloads map[string]any, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["expTime"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
