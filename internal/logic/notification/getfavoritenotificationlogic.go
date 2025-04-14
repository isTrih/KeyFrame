package notification

import (
	"context"
	"encoding/json"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFavoriteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteNotificationLogic {
	return &GetFavoriteNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoriteNotificationLogic) GetFavoriteNotification(req *types.GetNotificationRequest) (resp *types.GetNotificationMultiResponse, err error) {
	// 从JWT中获取用户ID
	userId := l.ctx.Value("userId").(json.Number)
	uid, _ := userId.Int64()

	// 获取收藏类型通知（分组后的），类型为3
	groupedItems, total, err := l.svcCtx.NotificationsModel.GetGroupedNotifications(
		l.ctx,
		uint64(uid),
		req.Page,
		req.Size,
		3, // 收藏类型为3
	)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	resp = &types.GetNotificationMultiResponse{
		Total:    total,
		MultiMSG: make([]types.MultiMSG, 0, len(groupedItems)),
	}

	for _, item := range groupedItems {
		resp.MultiMSG = append(resp.MultiMSG, types.MultiMSG{
			SenderNames:   item["sender_name"].([]string),
			SenderAvatars: item["sender_avatar"].([]string),
			Type:          int8(item["type"].(int64)),
			Total:         item["total"].(uint64),
			Time:          item["time"].(uint64),
		})
	}

	return resp, nil
}
