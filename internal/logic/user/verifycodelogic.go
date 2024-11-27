package user

import (
	"context"
	"math/rand"
	"strconv"
	"time"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewVerifyCodeLogic 获取验证码
func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCodeLogic) VerifyCode(req *types.VerifyCodeRequest) (resp *types.VerifyCodeResponse, err error) {
	// 生成验证码
	rand.New(rand.NewSource(time.Now().UnixNano())) // 取纳秒时间戳，可以保证每次的随机数种子都不同
	code := rand.Intn(900000) + 100000

	// 存储到redis中
	rdsErr := utils.RedisStorage(req.Mobile, code, 600)
	if rdsErr != nil {
		return nil, rdsErr
	}

	// 发送短信
	smsErr := utils.SendSms(req.Mobile, strconv.Itoa(code))
	if smsErr != nil {
		return nil, smsErr
	}

	resp = new(types.VerifyCodeResponse)
	resp.Status = "success"
	return resp, nil
}
