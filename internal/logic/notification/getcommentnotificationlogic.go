package notification

import (
	"context"
	"encoding/json"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/mdl/notifications"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentNotificationLogic {
	return &GetCommentNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentNotificationLogic) GetCommentNotification(req *types.GetNotificationRequest) (resp *types.GetNotificationResponse, err error) {
	// 从JWT中获取用户ID
	userId := l.ctx.Value("userId").(json.Number)
	uid, _ := userId.Int64()

	// 获取评论类型通知，类型为1
	items, total, err := l.svcCtx.NotificationsModel.GetNotificationsByType(
		l.ctx,
		uint64(uid),
		req.Page,
		req.Size,
		1, // 评论类型为1
	)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	resp = &types.GetNotificationResponse{
		Total:     total,
		SingleMSG: make([]types.SingleMSG, 0, len(items)),
	}

	for _, item := range items {
		var extra notifications.Extra
		_ = json.Unmarshal([]byte(item.Extra), &extra)

		resp.SingleMSG = append(resp.SingleMSG, types.SingleMSG{
			SenderId:      item.SenderId,
			SenderName:    extra.SenderName,
			SenderAvatar:  extra.SenderAvatar,
			Type:          int8(item.Type),
			Content:       item.Content,
			TargetContent: item.TargetContent,
			Time:          uint64(item.CreateTime.Unix()),
		})
	}

	return resp, nil
}
