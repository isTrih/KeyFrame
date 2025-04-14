package notification

import (
	"context"
	"encoding/json"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkAllNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllNotificationReadLogic {
	return &MarkAllNotificationReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAllNotificationReadLogic) MarkAllNotificationRead() (resp *types.MarkReadResponse, err error) {
	// 从JWT中获取用户ID
	userId := l.ctx.Value("userId").(json.Number)
	uid, _ := userId.Int64()

	// 标记所有通知为已读
	err = l.svcCtx.NotificationsModel.MarkAllAsRead(l.ctx, uint64(uid))
	if err != nil {
		return &types.MarkReadResponse{Success: false}, err
	}

	return &types.MarkReadResponse{Success: true}, nil
}
