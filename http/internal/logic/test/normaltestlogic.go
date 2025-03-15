package test

import (
	"context"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type NormalTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNormalTestLogic // 创建帧（文章）
func NewNormalTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NormalTestLogic {
	return &NormalTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NormalTestLogic) NormalTest(req *types.TestRequest) (resp *types.TestResponse, err error) {
	err = utils.UpdateFileLifecycle("img", req.Content, l.svcCtx.Config.Qiniu.AK, l.svcCtx.Config.Qiniu.SK, "chaozj-keyframe")
	if err != nil {
		return nil, err
	}
	//insp, err := utils.DoInsp(l.svcCtx.Config, req.Content)
	//resp = new(types.TestResponse)
	//resp.Content = strconv.Itoa(int(insp))
	//if insp != 0 {
	//	// 违规逻辑
	//}
	//if err != nil {
	//	return nil, err
	//}
	return
}
