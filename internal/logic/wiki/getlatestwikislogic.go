package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestWikisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLatestWikisLogic // 获取最新 Wiki 列表
func NewGetLatestWikisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestWikisLogic {
	return &GetLatestWikisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLatestWikisLogic) GetLatestWikis(req *types.GetLatestWikisRequest) (resp *types.WikiListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
