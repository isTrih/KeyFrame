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

const defaultEnd = 1 << 40

// FeedWithIndex 用于存储文章和对应的索引
type FeedWithIndex struct {
	Index int
	Feed  *article.Feeds
}

func (l *UserFeedsLogic) UserFeeds(req *types.UserFeedsRequest) (resp *types.GetIndexFeedsResponse, err error) {
	offset := req.Offset
	ftype := req.FeedType

	if offset != 0 && offset < 10 {
		resp = new(types.GetIndexFeedsResponse)
		resp.Status = "没有更多帖子"
		resp.Feeds = []types.Feed{}
		return resp, nil
	}

	if ftype == 2 {
		feedIds, _ := l.GetLikeIds(l.ctx, req.UserId, offset)
		if len(feedIds) > 0 {
			var tmp []uint64
			var feeds = []types.Feed{}
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
			}
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&Redis"
			resp.Feeds = feeds
			return resp, nil
		} else {
			var feeds = []types.Feed{}
			var tmp []uint64
			list, err := l.svcCtx.UserActionModel.GetUserLikeList(l.ctx, offset, req.UserId)
			println("list", list)
			if err != nil {
				return nil, err
			}
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
			}
			threading.GoSafe(func() {
				err = l.addCacheLike(context.Background(), list, req.UserId)
				if err != nil {
					logx.Error("addCacheLike failed: %v", err)
				}
			})
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&NO"
			resp.Feeds = feeds
			return resp, nil
		}
	} else if ftype == 1 {
		feedIds, _ := l.GetCollectIds(l.ctx, req.UserId, offset)
		if len(feedIds) > 0 {
			var tmp []uint64

			var feeds = []types.Feed{}
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
			}
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&Redis"
			resp.Feeds = feeds
			return resp, nil
		} else {
			var feeds = []types.Feed{}
			var tmp []uint64
			list, err := l.svcCtx.UserActionModel.GetUserCollectList(l.ctx, offset, req.UserId)
			if err != nil {
				return nil, err
			}
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
			}
			threading.GoSafe(func() {
				err = l.addCacheCollect(context.Background(), list, req.UserId)
				if err != nil {
					logx.Error("addCacheCollect failed: %v", err)
				}
			})
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&NO"
			resp.Feeds = feeds
			return resp, nil
		}

	} else {

		feedIds, _ := l.GetUploadIds(l.ctx, req.UserId, offset)
		if len(feedIds) > 0 {
			var tmp []uint64
			println("feedIds", feedIds)
			var feeds = []types.Feed{}
			list, err := l.GetUploadListByIds(l.ctx, feedIds)
			if err != nil {
				return nil, err
			}
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
				}
			}
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&Redis"
			resp.Feeds = feeds
			return resp, nil
		} else {
			var feeds = []types.Feed{}
			var tmp []uint64
			list, err := l.svcCtx.ArticleModel.GetUserUploadList(l.ctx, offset, req.UserId)
			if err != nil {
				return nil, err
			}
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
				}
			}
			threading.GoSafe(func() {
				err = l.addCacheUpload(context.Background(), list, req.UserId)
				if err != nil {
					logx.Error("addCacheUpload failed: %v", err)
				}
			})
			resp = new(types.GetIndexFeedsResponse)
			resp.Status = "Success&NO"
			resp.Feeds = feeds
			return resp, nil
		}
	}
}

func (l *UserFeedsLogic) GetCollectIds(ctx context.Context, uid, offset uint64) ([]uint64, error) {
	key := fmt.Sprintf("cache:keyframe:user:id:%d:collect", uid)

	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, defaultEnd, int(offset/10), 10)

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
	// 用于存储每个索引对应的文章
	result := make([]*article.Feeds, len(feedIds))

	_, err := mr.MapReduce[uint64, FeedWithIndex, struct{}](func(source chan<- uint64) {
		for _, fid := range feedIds {
			source <- fid
		}
	}, func(id uint64, writer mr.Writer[FeedWithIndex], cancel func(error)) {
		// 找到当前 id 在 feedIds 中的索引
		var index int
		for i, fid := range feedIds {
			if fid == id {
				index = i
				break
			}
		}

		p, err := l.svcCtx.ArticleModel.FindOneMix(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		// 写入带有索引的文章
		writer.Write(FeedWithIndex{
			Index: index,
			Feed:  p,
		})
	}, func(pipe <-chan FeedWithIndex, writer mr.Writer[struct{}], cancel func(error)) {
		// 从管道中读取带有索引的文章，并根据索引放入结果切片中
		for feedWithIndex := range pipe {
			result[feedWithIndex.Index] = feedWithIndex.Feed
		}
		// 这里不需要写入结果，因为已经填充到 result 切片中
		writer.Write(struct{}{})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (l *UserFeedsLogic) addCacheCollect(ctx context.Context, feedIds []*article.Feeds, userId uint64) error {

	if len(feedIds) == 0 {
		return nil
	}
	key := fmt.Sprintf("cache:keyframe:user:id:%d:collect", userId)
	for _, feedId := range feedIds {
		var score int64
		score = feedId.PublishTime.Unix()

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatUint(feedId.Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600)
}

func (l *UserFeedsLogic) GetLikeIds(ctx context.Context, uid, offset uint64) ([]uint64, error) {
	key := fmt.Sprintf("cache:keyframe:user:id:%d:like", uid)

	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, defaultEnd, int(offset/10), 10)

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
	// 用于存储每个索引对应的文章
	result := make([]*article.Feeds, len(feedIds))

	_, err := mr.MapReduce[uint64, FeedWithIndex, struct{}](func(source chan<- uint64) {
		for _, fid := range feedIds {
			source <- fid
		}
	}, func(id uint64, writer mr.Writer[FeedWithIndex], cancel func(error)) {
		// 找到当前 id 在 feedIds 中的索引
		var index int
		for i, fid := range feedIds {
			if fid == id {
				index = i
				break
			}
		}

		p, err := l.svcCtx.ArticleModel.FindOneMix(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		// 写入带有索引的文章
		writer.Write(FeedWithIndex{
			Index: index,
			Feed:  p,
		})
	}, func(pipe <-chan FeedWithIndex, writer mr.Writer[struct{}], cancel func(error)) {
		// 从管道中读取带有索引的文章，并根据索引放入结果切片中
		for feedWithIndex := range pipe {
			result[feedWithIndex.Index] = feedWithIndex.Feed
		}
		// 这里不需要写入结果，因为已经填充到 result 切片中
		writer.Write(struct{}{})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (l *UserFeedsLogic) addCacheLike(ctx context.Context, feedIds []*article.Feeds, userId uint64) error {

	if len(feedIds) == 0 {
		return nil
	}
	key := fmt.Sprintf("cache:keyframe:user:id:%d:like", userId)
	for _, feedId := range feedIds {
		var score int64
		score = feedId.PublishTime.Unix()

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatUint(feedId.Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, 3600)
}

func (l *UserFeedsLogic) GetUploadIds(ctx context.Context, uid, offset uint64) ([]uint64, error) {
	key := fmt.Sprintf("cache:keyframe:user:id:%d:upload", uid)
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, defaultEnd, int(offset/10), 10)
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
func (l *UserFeedsLogic) GetUploadListByIds(ctx context.Context, feedIds []uint64) ([]*article.Feeds, error) {
	// 用于存储每个索引对应的文章
	result := make([]*article.Feeds, len(feedIds))

	_, err := mr.MapReduce[uint64, FeedWithIndex, struct{}](func(source chan<- uint64) {
		for _, fid := range feedIds {
			source <- fid
		}
	}, func(id uint64, writer mr.Writer[FeedWithIndex], cancel func(error)) {
		// 找到当前 id 在 feedIds 中的索引
		var index int
		for i, fid := range feedIds {
			if fid == id {
				index = i
				break
			}
		}

		p, err := l.svcCtx.ArticleModel.FindOneMix(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		// 写入带有索引的文章
		writer.Write(FeedWithIndex{
			Index: index,
			Feed:  p,
		})
	}, func(pipe <-chan FeedWithIndex, writer mr.Writer[struct{}], cancel func(error)) {
		// 从管道中读取带有索引的文章，并根据索引放入结果切片中
		for feedWithIndex := range pipe {
			result[feedWithIndex.Index] = feedWithIndex.Feed
		}
		// 这里不需要写入结果，因为已经填充到 result 切片中
		writer.Write(struct{}{})
	})
	if err != nil {
		return nil, err
	}
	return result, nil

}
func (l *UserFeedsLogic) addCacheUpload(ctx context.Context, feedIds []*article.Feeds, userId uint64) error {

	if len(feedIds) == 0 {
		return nil
	}
	key := fmt.Sprintf("cache:keyframe:user:id:%d:upload", userId)
	for _, feedId := range feedIds {
		var score int64
		score = feedId.PublishTime.Unix()

		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatUint(feedId.Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, 480)

}
