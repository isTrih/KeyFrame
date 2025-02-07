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

type GetUpTokenChangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUpTokenChangeLogic // 获取到上传修改型token
func NewGetUpTokenChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpTokenChangeLogic {
	return &GetUpTokenChangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpTokenChangeLogic) GetUpTokenChange(req *types.ChangeResquest) (resp *types.UpResponse, err error) {
	accessKey := l.svcCtx.Config.Qiniu.AK
	secretKey := l.svcCtx.Config.Qiniu.SK
	mac := credentials.NewCredentials(accessKey, secretKey)
	bucket := "chaozj-keyframe"
	keyToOverwrite := req.Key //要覆盖的文件名
	putPolicy, err := uptoken.NewPutPolicyWithKey(bucket, keyToOverwrite, time.Now().Add(1*time.Hour))
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
