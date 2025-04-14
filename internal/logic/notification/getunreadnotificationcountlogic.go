package notification

import (
	"context"
	"encoding/json"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadNotificationCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadNotificationCountLogic {
	return &GetUnreadNotificationCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnreadNotificationCountLogic) GetUnreadNotificationCount() (resp *types.UnreadCountResponse, err error) {
	// 从JWT中获取用户ID
	userId := l.ctx.Value("userId").(json.Number)
	uid, _ := userId.Int64()

	// 获取各类型未读通知数
	commentCount, err := l.svcCtx.NotificationsModel.GetTypeUnreadCount(l.ctx, uint64(uid), 1) // 评论
	if err != nil {
		return nil, err
	}

	likeCount, err := l.svcCtx.NotificationsModel.GetTypeUnreadCount(l.ctx, uint64(uid), 2) // 点赞
	if err != nil {
		return nil, err
	}

	favoriteCount, err := l.svcCtx.NotificationsModel.GetTypeUnreadCount(l.ctx, uint64(uid), 3) // 收藏
	if err != nil {
		return nil, err
	}

	followCount, err := l.svcCtx.NotificationsModel.GetTypeUnreadCount(l.ctx, uint64(uid), 4) // 关注
	if err != nil {
		return nil, err
	}

	// 获取总未读数
	totalCount, err := l.svcCtx.NotificationsModel.GetUnreadCount(l.ctx, uint64(uid))
	if err != nil {
		return nil, err
	}

	resp = &types.UnreadCountResponse{
		CommentCount:  uint64(commentCount),
		LikeCount:     uint64(likeCount),
		FavoriteCount: uint64(favoriteCount),
		FollowCount:   uint64(followCount),
		TotalCount:    uint64(totalCount),
	}

	return resp, nil
}
