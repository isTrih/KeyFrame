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

// NewGetUpTokenLogic 获取到上传token
func NewGetUpTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpTokenLogic {
	return &GetUpTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpTokenLogic) GetUpToken(req *types.GetUpTokenRequest) (resp *types.GetUpTokenResponse, err error) {
	accessKey := l.svcCtx.Config.Qiniu.AK
	secretKey := l.svcCtx.Config.Qiniu.SK
	mac := credentials.NewCredentials(accessKey, secretKey)
	bucket := "chaozj-keyframe"
	// 上传凭证有效期5分钟
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}
	// 上传文件的key为文件类型+文件hash
	putPolicy.SetSaveKey(req.Type + "/$(etag)")
	putPolicy.SetForceSaveKey(true)

	putPolicy.SetReturnBody(`{"key": $(etag),"w": $(imageInfo.width), "h": $(imageInfo.height)}`)
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return nil, err
	}
	resp = new(types.GetUpTokenResponse)
	resp.Token = upToken
	return resp, nil
}
