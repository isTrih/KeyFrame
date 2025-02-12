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
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()
	user, _ := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	encryptedPassword := EncryptPassword(req.NewPassword)
	if encryptedPassword != user.Password {
		// 可以更改密码

		// 更改密码
		err = l.svcCtx.UserModel.UpdatePassword(l.ctx, uint64(uid), encryptedPassword)
		resp = &types.ChangePasswordResponse{
			Status: "SUCCESS",
		}
		return resp, nil
	}
	return nil, errors.New(4003, "新密码与旧密码相同")
}
