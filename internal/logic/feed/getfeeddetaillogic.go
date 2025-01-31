package feed

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/x/errors"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetFeedDetailLogic // 帧（文章）详情
func NewGetFeedDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedDetailLogic {
	return &GetFeedDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedDetailLogic) GetFeedDetail(req *types.GetFeedDetailRequest) (resp *types.GetFeedDetailResponse, err error) {
	// 查询文章
	a, err := l.svcCtx.ArticleModel.FindOneDetail(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var urls []string
	err = json.Unmarshal([]byte(a.MediaList.String), &urls)
	if err != nil {
		return nil, errors.New(4003, "解析JSON失败:"+err.Error())
	}

	resp = &types.GetFeedDetailResponse{
		Feed: types.FeedDetail{
			Title:      a.Title,
			Id:         a.Id,
			MediaUrl:   a.CoverUrl,
			Content:    a.Content,
			Type:       a.Type,
			CommentNum: a.CommentNum,
			LikeNum:    a.LikeNum,
			CollectNum: a.CollectNum,
			ViewNum:    a.ViewNum,
			ShareNum:   a.ShareNum,
			MediaInfo: types.MediaInfo{
				Width:  a.Width,
				Height: a.Height,
			},
			MediaList: urls,
			Author: types.FeedUser{
				Id:       a.AuthorId,
				UserName: a.UserName,
				Avatar:   a.Avatar,
			},
			PublishTime: uint64(a.PublishTime.Unix()),
			CreateTime:  uint64(a.CreateTime.Unix()),
			UpdateTime:  uint64(a.UpdateTime.Unix()),
			AiInsp:      a.AiInsp,
			AiInspCode:  a.AiInspCode,
			Insp:        a.Insp,
		},
	}

	return resp, nil
}
