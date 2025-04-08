package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/x/errors"
	"time"
	"zerobackend/internal/utils"
	model "zerobackend/mdl/user"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	insp, err := utils.DoInsp(l.svcCtx.Config, req.Username)
	if insp != 0 {
		// 违规逻辑
		return nil, errors.New(6099, "昵称有违规内容")
	}

	var user sql.Result

	rds, rds2 := utils.RedisCheck(l.svcCtx.Config, req.Mobile, req.VerifyCode)
	if rds != nil {
		return nil, rds
	}
	if rds2 != nil {
		return nil, rds2
	}
	focusId := utils.DecryptTriDESToNumber(l.svcCtx.Config.InviteKey.KEY,
		l.svcCtx.Config.InviteKey.IV, req.CZJCode)
	if focusId != 0 {
		check, checkerr := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
		if checkerr != nil && checkerr != model.ErrNotFound {
			fmt.Println(checkerr)
			return nil, errors.New(4003, "查询数据失败")
		}
		if check == nil {
			user, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
				Id:        focusId,
				Mobile:    req.Mobile,
				Nickname:  req.Username,
				Password:  EncryptPassword(req.Password),
				Signature: "",
				Avatar:    "avatar.jpg",
			})
			if err != nil {
				fmt.Println(err)
				return nil, errors.New(6011, "注册失败")
			}
			logc.Info(l.ctx, user)
		} else {
			return nil, errors.New(6012, "邀请码已使用")
		}
	}

	check, checkerr := l.svcCtx.UserModel.FindOneByMobile(l.ctx, req.Mobile)
	if checkerr != nil && checkerr != model.ErrNotFound {
		fmt.Println(checkerr)
		return nil, errors.New(4003, "查询数据失败")
	}
	if check == nil {
		user, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
			Mobile:    req.Mobile,
			Nickname:  req.Username,
			Password:  EncryptPassword(req.Password),
			Signature: "",
			Avatar:    "avatar.jpg",
		})
		if err != nil {
			fmt.Println(err)
			return nil, errors.New(6011, "注册失败")
		}
		logc.Info(l.ctx, user)
	} else {
		return nil, errors.New(6012, "用户已存在")
	}
	payloads := make(map[string]any)
	uid, _ := user.LastInsertId()
	uintUid := uint64(uid)
	payloads["UID"], _ = user.LastInsertId()
	payloads["UTYPE"] = 0
	payloads["USTATUS"] = 0

	accessToken, tokenErr := utils.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}
	resp = new(types.RegisterResponse)
	resp.UserId = uintUid
	resp.UserName = req.Username
	resp.Avatar = "avatar.jpg"
	resp.Signature = ""
	resp.Token = accessToken
	resp.UserType = 0
	return resp, nil
}

// EncryptPassword 加密密码
// 示例:
//
//	encryptedPassword := EncryptPassword("123456")
func EncryptPassword(needEncryptPassword string) (encryptedPassword string) {

	encryptedPassword = fmt.Sprintf("%x", sha256.Sum256([]byte(needEncryptPassword+"tmh")))
	return encryptedPassword
}
