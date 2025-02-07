package upload

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"time"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUpTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUpTokenLogic // 获取到上传token
func NewGetUpTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpTokenLogic {
	return &GetUpTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpTokenLogic) GetUpToken() (resp *types.UpResponse, err error) {
	accessKey := l.svcCtx.Config.Qiniu.AK
	secretKey := l.svcCtx.Config.Qiniu.SK
	mac := credentials.NewCredentials(accessKey, secretKey)
	bucket := "chaozj-keyframe"
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return nil, err
	}
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return nil, err
	}
	resp = new(types.UpResponse)
	resp.Token = upToken
	return resp, nil
}
