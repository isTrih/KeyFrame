package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWikiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateWikiLogic // 编辑 Wiki
func NewUpdateWikiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWikiLogic {
	return &UpdateWikiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWikiLogic) UpdateWiki(req *types.UpdateWikiRequest) (resp *types.StatusResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
