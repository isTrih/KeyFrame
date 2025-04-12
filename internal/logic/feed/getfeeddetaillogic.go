package feed

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/x/errors"
	"strconv"
	"zerobackend/internal/nats/producer"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/mdl/article_metrics"

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

	articleMetric, amErr := l.svcCtx.ArticleMetricsModel.FindOneByArticleId(l.ctx, int64(req.Id))
	switch amErr {
	case nil:
		break
	case sqlc.ErrNotFound:
		articleMetric = &article_metrics.ArticleMetrics{
			Likes:    0,
			Collects: 0,
			Comments: 0,
		}
	default:
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
			CommentNum: uint64(articleMetric.Comments),
			LikeNum:    uint64(articleMetric.Likes),
			CollectNum: uint64(articleMetric.Collects),
			MediaInfo: types.MediaInfo{
				Width:  a.Width,
				Height: a.Height,
			},
			MediaList: urls,
			Author: types.FeedUser{
				Id:       a.AuthorId,
				UserName: a.UserName,
				Avatar:   a.Avatar,
				Type:     a.Type,
			},
			PublishTime: uint64(a.PublishTime.Unix()),
			IpLocation:  a.IpLocation,
		},
	}

	Queeerr := producer.SendMessageToQueue(strconv.FormatInt(int64(timex.Now()), 10))
	if Queeerr != nil {
		return nil, Queeerr
	}

	return resp, nil
}
