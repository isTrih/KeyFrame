package feed

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/x/errors"
	"strconv"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLikeFeedLogic // 点赞帧（文章）
func NewLikeFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeFeedLogic {
	return &LikeFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeFeedLogic) LikeFeed(req *types.DeleteFeedRequest) (resp *types.StatusResponse, err error) {
	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()
	// 计数器
	store, _ := redis.NewRedis(l.svcCtx.Config.BizRedis)
	limiter := limit.NewPeriodLimit(30, 60, store, "action")

	timer, _ := limiter.Take(strconv.FormatInt(uid, 10))
	switch timer {
	case limit.Allowed:
		// todo: 添加逻辑
		err = l.svcCtx.ArticleModel.LikeArticle(l.ctx, uid, int64(req.Id))
		if err != nil {
			return nil, err
		}
		// 构造返回值
		resp = &types.StatusResponse{
			Status: "ok",
		}
		return resp, nil
	case limit.HitQuota:
		return nil, errors.New(1002, "操作过于频繁，请30s后再试")
	case limit.OverQuota:
		return nil, errors.New(1002, "操作过于频繁，请30后再试")
	default:
		return nil, errors.New(1003, "发生未知错误")
	}
}
