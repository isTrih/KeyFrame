package user

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/x/errors"
	"zerobackend/internal/utils"
	model "zerobackend/mdl/user"

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
	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()

	insp, err := utils.DoInsp(l.svcCtx.Config, req.Nickname+req.Signature)
	if insp != 0 {
		// 违规逻辑
		return nil, errors.New(6099, "昵称或签名有违规内容")
	}
	oldUser, _ := l.svcCtx.UserModel.FindOne(l.ctx, uid)

	err = l.svcCtx.UserModel.Update(l.ctx, &model.User{
		Id:         uid,
		Nickname:   req.Nickname,
		Avatar:     req.Avatar,
		Signature:  req.Signature,
		Type:       oldUser.Type,
		Vnote:      oldUser.Vnote,
		Mobile:     oldUser.Mobile,
		Status:     oldUser.Status,
		BannedTime: oldUser.BannedTime,
		BanTime:    oldUser.BanTime,
		IpAddress:  oldUser.IpAddress,
		IpLocation: oldUser.IpLocation,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ChangeUserInfoResponse{
		Status: "ok",
	}
	return resp, nil
}
