package feed

import (
	"context"
	"encoding/json"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedListSelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetFeedListSelfLogic // 创作中心获取帧（文章）
func NewGetFeedListSelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedListSelfLogic {
	return &GetFeedListSelfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedListSelfLogic) GetFeedListSelf(req *types.GetFeedListRequest) (resp *types.GetFeedListResponse, err error) {
	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()

	self, err := l.svcCtx.ArticleModel.GetUserUploadListForSelf(l.ctx, req.Offset, uint64(uid))
	if err != nil {
		return nil, err
	}

	feedCount, _ := l.svcCtx.ArticleModel.GetFeedsNum(l.ctx, uint64(uid))

	var convertedFeeds []types.Feeds
	for _, feed := range self {
		convertedFeeds = append(convertedFeeds, types.Feeds{
			Title:       feed.Title,
			Id:          feed.Id,
			AuthorId:    feed.AuthorId,
			UserName:    feed.UserName,
			Avatar:      feed.Avatar,
			CoverUrl:    feed.CoverUrl,
			Height:      feed.Height,
			Width:       feed.Width,
			LikeNum:     feed.LikeNum,
			CollectNum:  feed.CollectNum,
			CommentNum:  feed.CommentNum,
			AiInsp:      feed.AiInsp,
			Insp:        feed.Insp,
			PublishTime: uint64(feed.PublishTime.Unix()),
		})
	}

	resp = &types.GetFeedListResponse{
		Feeds:  convertedFeeds,
		Status: "ok",
		Total:  uint64(feedCount),
	}
	return
}
