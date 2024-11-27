package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/x/errors"
	"strconv"
	"time"
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
	conf := redis.RedisConf{
		Host:        "106.54.6.216:6379",
		Type:        "node",
		Pass:        "chaozj123123.",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	code, err := rds.Get(req.Mobile)
	if err != nil {
		logc.Error(l.ctx, err)
	}

	if code != req.VerifyCode && code != "" {
		fmt.Println(code)
		return nil, errors.New(1001, "验证码错误")
	}
	if code == "" {
		fmt.Println(code)
		return nil, errors.New(1002, "验证码过期或未获取")
	}

	check, checkerr := l.svcCtx.UserModel.FindOneBymobile(l.ctx, mobileInt)
	if checkerr != nil && checkerr != model.ErrNotFound {
		fmt.Println(checkerr)
		return nil, errors.New(4001, "查询数据失败")
	}
	if check == nil {
		psw := sha256.Sum256([]byte(req.Password + "tmh"))
		user, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
			Mobile:    req.Mobile,
			Nickname:  req.Username,
			Password:  fmt.Sprintf("%x", psw),
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

	accessToken, tokenErr := l.GetToken(time.Now().Unix(), l.svcCtx.Config.Auth.AccessSecret, payloads, l.svcCtx.Config.Auth.AccessExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}
	resp = new(types.RegisterResponse)
	resp.UserId, _ = user.LastInsertId()
	resp.Token = accessToken

	rds.Del(req.Mobile)
	return resp, nil
}

func (l *RegisterLogic) GetToken(iat int64, secretKey string, payloads map[string]any, seconds int64) (string, error) {
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
