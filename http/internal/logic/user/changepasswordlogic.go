package user

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/x/errors"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangePasswordLogic 更改密码 需要token
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (resp *types.ChangePasswordResponse, err error) {
	if req.NewPassword != req.OldPassword {
		// 可以更改密码
		uidjson, _ := l.ctx.Value("UID").(json.Number)
		uid, _ := uidjson.Int64()
		// 更改密码
		err = l.svcCtx.UserModel.UpdatePassword(l.ctx, uint64(uid), EncryptPassword(req.NewPassword))
		resp = &types.ChangePasswordResponse{
			Status: "SUCCESS",
		}
		return resp, nil
	}
	return nil, errors.New(4003, "新密码与旧密码相同")
}
