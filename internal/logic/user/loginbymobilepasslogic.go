package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/x/errors"
	"net"
	"time"
	"zerobackend/internal/utils"
	model "zerobackend/mdl/user"

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
	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4003, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")

	}

	// 校验密码
	if user.Password != EncryptPassword(req.Password) {
		return nil, errors.New(6032, "密码错误")
	}

	// 查询ip
	var ip = req.KIP
	region, err := l.svcCtx.IP4Searcher.Search(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}
	fmt.Println(region, net.ParseIP(ip))
	//更新IP信息以及归属地
	upErr := l.svcCtx.UserModel.UpdateIpByMobile(l.ctx, req.Mobile, region, ip)
	if upErr != nil {
		fmt.Println(upErr)
		return nil, errors.New(4004, "更新数据失败，请联系管理员")
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

	//返回正确的token
	resp = new(types.LoginResponse)
	resp.Token = accessToken
	resp.UserType = uint8(user.Type)
	resp.Avatar = user.Avatar
	resp.Signature = user.Signature
	//resp.Signature = req.XRI
	resp.UserName = user.Nickname
	resp.UserId = user.Id

	return resp, nil
}
