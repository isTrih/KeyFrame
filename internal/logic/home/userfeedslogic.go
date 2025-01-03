package home

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"slices"
	"strconv"
	"zerobackend/mdl/article"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFeedsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserFeedsLogic // 获取用户主页帖子
func NewUserFeedsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFeedsLogic {
	return &UserFeedsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFeedsLogic) UserFeeds(req *types.UserFeedsRequest) (resp *types.GetIndexFeedsResponse, err error) {
	offset := req.Offset
	ftype := req.FeedType

	if offset != 0 && offset < 10 {
		return nil, nil
	}

	if ftype == "喜欢" {
		feedIds, _ := l.GetLikeIds(l.ctx, req.UserId, offset)
		if len(feedIds) > 0 {
			var tmp []uint64

			var feeds []types.Feed
			list, err := l.GetLikeListByIds(l.ctx, feedIds)
			if err != nil {
				return nil, err
			}
			for _, v := range list {
				if slices.Contains(tmp, v.Id) == false {
					tmp = append(tmp, v.Id)
					feeds = append(feeds, types.Feed{
						Id:      v.Id,
						Title:   v.Title,
						ViewNum: v.Views,
						LikeNum: v.Likes,
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
			}
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "success+redis"
			resp.Feeds = feeds
			return resp, nil
		} else {
			var feeds []types.Feed
			var tmp []uint64
			list, err := l.svcCtx.LikeRecordModel.GetUserLikeList(l.ctx, offset, req.UserId)
			if err != nil {
				return nil, err
			}
			for _, v := range list {
				if slices.Contains(tmp, v.Id) == false {
					tmp = append(tmp, v.Id)
					feeds = append(feeds, types.Feed{
						Id:      v.Id,
						Title:   v.Title,
						ViewNum: v.Views,
						LikeNum: v.Likes,
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
			}
			threading.GoSafe(func() {
				err = l.addCacheLike(context.Background(), list, req.UserId)
				if err != nil {
					logx.Error("addCacheLike failed: %v", err)
				}
			})
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "success+noRedis"
			resp.Feeds = feeds
			return resp, nil
		}
	} else if ftype == "收藏" {
		feedIds, _ := l.GetCollectIds(l.ctx, req.UserId, offset)
		if len(feedIds) > 0 {
			var tmp []uint64

			var feeds []types.Feed
			list, err := l.GetCollectListByIds(l.ctx, feedIds)
			if err != nil {
				return nil, err
			}
			for _, v := range list {
				if slices.Contains(tmp, v.Id) == false {
					tmp = append(tmp, v.Id)
					feeds = append(feeds, types.Feed{
						Id:      v.Id,
						Title:   v.Title,
						ViewNum: v.Views,
						LikeNum: v.Likes,
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
			}
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "success"
			resp.Feeds = feeds
			return resp, nil
		} else {
			var feeds []types.Feed
			var tmp []uint64
			list, err := l.svcCtx.CollectModel.GetUserCollectList(l.ctx, offset, req.UserId)
			if err != nil {
				return nil, err
			}
			for _, v := range list {
				if slices.Contains(tmp, v.Id) == false {
					tmp = append(tmp, v.Id)
					feeds = append(feeds, types.Feed{
						Id:      v.Id,
						Title:   v.Title,
						ViewNum: v.Views,
						LikeNum: v.Likes,
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
			}
			threading.GoSafe(func() {
				err = l.addCacheCollect(context.Background(), list, req.UserId)
				if err != nil {
					logx.Error("addCacheCollect failed: %v", err)
				}
			})
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "success"
			resp.Feeds = feeds
			return resp, nil
		}

	} else {
		//TODO:有空重构一下这里的代码
		list, err := l.svcCtx.ArticleModel.GetUserFeeds(l.ctx, offset, req.UserId)
		if err != nil {
			return nil, err
		}
		fmt.Println(list)
		var feeds []types.Feed
		var tmp []uint64
		for _, v := range list {
			if slices.Contains(tmp, v.Id) == false {
				tmp = append(tmp, v.Id)
				feeds = append(feeds, types.Feed{
					Id:      v.Id,
					Title:   v.Title,
					ViewNum: v.Views,
					LikeNum: v.Likes,
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
}

func (l *UserFeedsLogic) GetCollectIds(ctx context.Context, uid, offset uint64) ([]uint64, error) {
	key := fmt.Sprintf("user:collect:id:%d", uid)

	pairs, err := l.svcCtx.BizRedis.ZrangebyscoreWithScoresAndLimitCtx(ctx, key, int64(offset), int64(offset+10), 0, 10)

	if err != nil {
		return nil, err
	}

	var ids []uint64
	for _, pair := range pairs {
		id, err := strconv.ParseUint(pair.Key, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func (l *UserFeedsLogic) GetCollectListByIds(ctx context.Context, feedIds []uint64) ([]*article.Feeds, error) {
	feeds, err := mr.MapReduce[uint64, *article.Feeds, []*article.Feeds](func(source chan<- uint64) {
		for _, fid := range feedIds {
			source <- fid
		}
	}, func(id uint64, writer mr.Writer[*article.Feeds], cancel func(error)) {
		p, err := l.svcCtx.ArticleModel.FindOneMix(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *article.Feeds, writer mr.Writer[[]*article.Feeds], cancel func(error)) {
		var feeds []*article.Feeds
		for feed := range pipe {
			feeds = append(feeds, feed)
		}
		writer.Write(feeds)
	})
	if err != nil {
		return nil, err
	}
	return feeds, nil
}
func (l *UserFeedsLogic) addCacheCollect(ctx context.Context, feedIds []*article.Feeds, userId uint64) error {

	if len(feedIds) == 0 {
		return nil
	}
	key := fmt.Sprintf("user:collect:id:%d", userId)
	for _, feedId := range feedIds {
		var score int64
		score = int64(feedId.Id)

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatUint(feedId.Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600)
}

func (l *UserFeedsLogic) GetLikeIds(ctx context.Context, uid, offset uint64) ([]uint64, error) {
	key := fmt.Sprintf("user:like:id:%d", uid)

	pairs, err := l.svcCtx.BizRedis.ZrangebyscoreWithScoresAndLimitCtx(ctx, key, int64(offset), int64(offset+10), 0, 10)

	if err != nil {
		return nil, err
	}

	var ids []uint64
	for _, pair := range pairs {
		id, err := strconv.ParseUint(pair.Key, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func (l *UserFeedsLogic) GetLikeListByIds(ctx context.Context, feedIds []uint64) ([]*article.Feeds, error) {
	feeds, err := mr.MapReduce[uint64, *article.Feeds, []*article.Feeds](func(source chan<- uint64) {
		for _, fid := range feedIds {
			source <- fid
		}
	}, func(id uint64, writer mr.Writer[*article.Feeds], cancel func(error)) {
		p, err := l.svcCtx.ArticleModel.FindOneMix(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *article.Feeds, writer mr.Writer[[]*article.Feeds], cancel func(error)) {
		var feeds []*article.Feeds
		for feed := range pipe {
			feeds = append(feeds, feed)
		}
		writer.Write(feeds)
	})
	if err != nil {
		return nil, err
	}
	return feeds, nil
}
func (l *UserFeedsLogic) addCacheLike(ctx context.Context, feedIds []*article.Feeds, userId uint64) error {

	if len(feedIds) == 0 {
		return nil
	}
	key := fmt.Sprintf("user:like:id:%d", userId)
	for _, feedId := range feedIds {
		var score int64
		score = int64(feedId.Id)

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatUint(feedId.Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600)
}
