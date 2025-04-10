package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/x/errors"
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
	rdsErr := utils.RedisStorage(l.svcCtx.Config, req.Mobile, code, 600)
	if rdsErr != nil {
		return nil, rdsErr
	}

	store, _ := redis.NewRedis(l.svcCtx.Config.BizRedis)
	// 计数器
	limiter := limit.NewPeriodLimit(60, 3, store, "smsCode")

	timer, _ := limiter.Take(req.Mobile)
	switch timer {
	case limit.Allowed:
		// 发送短信
		_, smsErr := utils.SendSms(l.svcCtx.Config, req.Mobile, strconv.Itoa(code))
		if smsErr != nil {
			resp = new(types.VerifyCodeResponse)
			resp.Status = "temp"
			resp.TempCode = strconv.Itoa(code)
			return resp, nil
		}
		resp = new(types.VerifyCodeResponse)
		resp.Status = "success"
		resp.TempCode = ""
		return resp, nil
	case limit.HitQuota:
		resp = new(types.VerifyCodeResponse)
		resp.Status = "HitQuota"
		resp.TempCode = strconv.Itoa(code)
		return resp, nil
	case limit.OverQuota:
		resp = new(types.VerifyCodeResponse)
		resp.Status = "OverQuota"
		resp.TempCode = strconv.Itoa(code)
		return resp, nil
	default:
		return nil, errors.New(1003, "发生未知错误")
	}

}
