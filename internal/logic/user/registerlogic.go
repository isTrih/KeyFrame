package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/x/errors"
	"strconv"
	"time"
	"zerobackend/internal/utils"
	"zerobackend/mdl/user/model"

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

	mobileInt, _ := strconv.ParseUint(req.Mobile, 10, 64)
	var user sql.Result

	rds, rds2 := utils.RedisCheck(req.Mobile, req.VerifyCode)
	if rds != nil {
		return nil, rds
	}
	if rds2 != nil {
		return nil, rds2
	}

	check, checkerr := l.svcCtx.UserModel.FindOneBymobile(l.ctx, mobileInt)
	if checkerr != nil && checkerr != model.ErrNotFound {
		fmt.Println(checkerr)
		return nil, errors.New(4001, "查询数据失败")
	}
	if check == nil {
		user, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
			Mobile:    req.Mobile,
			Nickname:  req.Username,
			Password:  EncryptPassword(req.Password),
			Signature: "CHAOZJ",
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
	payloads["UID"], _ = user.LastInsertId()
	payloads["UTYPE"] = 0

	accessToken, tokenErr := utils.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}
	resp = new(types.RegisterResponse)
	resp.UserId, _ = user.LastInsertId()
	resp.Token = accessToken
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
