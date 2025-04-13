package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWikiDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetWikiDetailLogic // 获取 Wiki 详情
func NewGetWikiDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWikiDetailLogic {
	return &GetWikiDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWikiDetailLogic) GetWikiDetail(req *types.GetWikiDetailRequest) (resp *types.GetWikiDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
