package article

import (
	"context"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNewArticleLogic 创建帧（文章）
func NewNewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewArticleLogic {
	return &NewArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewArticleLogic) NewArticle(req *types.NewArticleRequest) (resp *types.StatusResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
