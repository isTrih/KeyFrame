package wiki

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchWikiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSearchWikiLogic // 搜索 Wiki
func NewSearchWikiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchWikiLogic {
	return &SearchWikiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchWikiLogic) SearchWiki(req *types.SearchWikiRequest) (resp *types.WikiListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
