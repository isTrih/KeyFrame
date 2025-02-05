package home

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/url"
	"slices"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
)

type GetIndexFeedsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetIndexFeedsLogic 获取首页信息流
func NewGetIndexFeedsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIndexFeedsLogic {
	return &GetIndexFeedsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIndexFeedsLogic) GetIndexFeeds(req *types.GetIndexFeedsRequest) (resp *types.GetIndexFeedsResponse, err error) {
	offset := req.Offset

	query, err := url.QueryUnescape(req.Query) // 解码
	if err != nil {
		fmt.Println("Decoded Error:", err)
		return nil, err
	}
	list, err := l.svcCtx.ArticleModel.GetFeeds(l.ctx, offset, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(list)
	var feeds = []types.Feed{}
	var tmp []uint64
	for _, v := range list {
		if slices.Contains(tmp, v.Id) == false {
			tmp = append(tmp, v.Id)
			feeds = append(feeds, types.Feed{
				Id:      v.Id,
				Title:   v.Title,
				LikeNum: v.LikeNum,
				User: types.FeedUser{
					Id:       v.AuthorId,
					UserName: v.UserName,
					Avatar:   v.Avatar,
				},
				MediaUrl: v.CoverUrl,
				MediaInfo: types.MediaInfo{
					Height: v.Height,
					Width:  v.Width,
				},
				Loaded: false,
			})
			fmt.Println(v)
		}
		//var media types.FeedMedia
		//var user types.FeedUser
		//err := json.Unmarshal([]byte(v.Media), &media)
		//err = json.Unmarshal([]byte(v.User), &user)
		//if err != nil {
		//	return nil, err
		//}

		//feeds = append(feeds, types.Feed{
		//	Id:     v.Id,
		//	Loaded: false,
		//	Title:  v.Title,
		//	User:   user,
		//})
	}
	resp = new(types.GetIndexFeedsResponse)
	resp.Status = "success"
	resp.Feeds = feeds
	return resp, nil
}
