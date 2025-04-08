package user

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserRelationLogic // 获取用户关注\收藏\点赞列表
func NewUserRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRelationLogic {
	return &UserRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRelationLogic) UserRelation() (resp *types.UserRelationResponse, err error) {
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()
	action, err := l.svcCtx.UserActionModel.GetUserActionHistory(l.ctx, uint64(uid))
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(4003, err.Error())
	}

	// 构造返回值
	resp = &types.UserRelationResponse{
		LikeCommentList: action.LikeCommentList,
		LikeFeedList:    action.LikeFeedList,
		CollectFeedList: action.CollectFeedList,
		FollowList:      action.FollowList,
	}
	return resp, nil
}
