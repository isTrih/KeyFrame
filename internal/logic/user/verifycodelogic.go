package user

import (
	"context"
	"fmt"
	unisms "github.com/apistd/uni-go-sdk/sms"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math/rand"
	"strconv"
	"time"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取验证码
func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCodeLogic) VerifyCode(req *types.VerifyCodeRequest) (resp *types.VerifyCodeResponse, err error) {
	// 生成随机数
	rand.New(rand.NewSource(time.Now().UnixNano())) // 取纳秒时间戳，可以保证每次的随机数种子都不同
	code := rand.Intn(900000) + 100000

	conf := redis.RedisConf{
		Host:        "106.54.6.216:6379",
		Type:        "node",
		Pass:        "chaozj123123.",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	err = rds.Setex(req.Mobile, strconv.Itoa(code), 600)
	if err != nil {
		logc.Error(l.ctx, err)
	}

	// 初始化
	client := unisms.NewClient("n6Z4ZVVUToSXF8ktQbXeNfGvQRp6mcdoH43pjNhy3uKPXMCW8") // 若使用简易验签模式仅传入第一个参数即可

	// 构建信息
	message := unisms.BuildMessage()
	message.SetTo(req.Mobile)
	message.SetSignature("超正经科技")
	message.SetTemplateId("pub_verif_ttl3")
	message.SetTemplateData(map[string]string{"code": strconv.Itoa(code), "ttl": "10"}) // 设置自定义参数 (变量短信)

	// 发送短信
	res, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	resp = new(types.VerifyCodeResponse)
	resp.Status = "success"
	return resp, nil
}
