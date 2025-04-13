package feed

import (
	"context"
	"encoding/json"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteFeedLogic // 删除帧（文章）
func NewDeleteFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFeedLogic {
	return &DeleteFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFeedLogic) DeleteFeed(req *types.DeleteFeedRequest) (resp *types.StatusResponse, err error) {
	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()
	// 软删除
	err = l.svcCtx.ArticleModel.DeleteFeed(l.ctx, uid, int64(req.Id))
	if err != nil {
		return nil, err
	}
	resp = &types.StatusResponse{
		Status: "ok",
	}
	return resp, nil
}
