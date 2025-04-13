package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
	"zerobackend/internal/nats/producer"
	"zerobackend/internal/types"
)

var _ NotificationsModel = (*customNotificationsModel)(nil)

type (
	// NotificationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotificationsModel.
	NotificationsModel interface {
		notificationsModel
		CheckNotification(ctx context.Context, uid, receiverid uint64, Type, TargetType uint8) (bool, error)
	}

	customNotificationsModel struct {
		*defaultNotificationsModel
	}

	Extra struct {
		SenderAvatar string `json:"sender_avatar"`
		SenderName   string `json:"sender_name"`
	}
	NotiItems struct {
		Id            int64          `db:"id"`             // 主键，自增 ID
		SenderId      int64  `db:"sender_id"`      // 发送人 ID，系统通知可为 NULL
		ReceiverId    int64          `db:"receiver_id"`    // 接收人 ID
		Type          int64          `db:"type"`           // 通知类型，例如 1点赞2评论3收藏4关注
		Content       string `db:"content"`        // 通知显示内容，可用于消息列表展示
		TargetId      int64  `db:"target_id"`      // 通知对象的 ID
		IsRead        bool           `db:"is_read"`        // 是否已读标志
		CreateTime    time.Time      `db:"create_time"`    // 通知创建时间
		Extra         string         `db:"extra"`          // 额外信息（JSON 格式），如帖子标题、评论摘要等
		TargetContent string `db:"target_content"` // 目标内容
		TargetType    int64          `db:"target_type"`    // 目标类型 文章，评论，用户
	}
)

// NewNotificationsModel returns a model for the database table.
func NewNotificationsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotificationsModel {
	return &customNotificationsModel{
		defaultNotificationsModel: newNotificationsModel(conn, c, opts...),
	}
}

func (m *defaultNotificationsModel) InsertNotifications(ctx context.Context, ReceiverId, SenderId uint64, Type uint8, Content, TargetContent, SenderName, SenderAvatar string) error {

	msgTo := fmt.Sprintf("KEYFRAME.MSG.%d", ReceiverId)

	msg := &types.SingleMSG{
		SenderId:      int64(SenderId),
		SenderName:    SenderName,
		SenderAvatar:  SenderAvatar,
		Type:          int8(Type),
		Content:       Content,
		TargetContent: TargetContent,
		Time:          uint64(time.Now().Unix()),
	}
	msgJSON, _ := json.Marshal(msg)
	Queeerr := producer.SendMessageToQueue(msgTo, string(msgJSON))
	if Queeerr != nil {
		return Queeerr
	}

	extra := &Extra{
		SenderAvatar: SenderAvatar,
		SenderName:   SenderName,
	}
	extraJson, _ := json.Marshal(extra)

	_, err := m.Insert(ctx, &Notifications{
		SenderId:      sql.NullInt64{Int64: int64(SenderId), Valid: true},
		ReceiverId:    int64(ReceiverId),
		Type:          int64(Type),
		Content:       sql.NullString{String: Content, Valid: true},
		TargetId:      sql.NullInt64{},
		IsRead:        false,
		TargetContent: sql.NullString{String: TargetContent, Valid: true},
		Extra:         string(extraJson),
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *defaultNotificationsModel) CheckNotification(ctx context.Context, uid, receiverid uint64, Type, TargetType uint8) (bool, error) {
	var exists bool
	query := `
	WITH notification_check AS (
		SELECT 1 
		FROM "public"."notifications" 
		WHERE sender_id = $1 
		AND receiver_id = $2 
		AND type = $3 
		AND target_type = $4
		LIMIT 1
	)
	SELECT EXISTS (SELECT 1 FROM notification_check)`
	
	err := m.QueryRowNoCacheCtx(ctx, &exists, query, 
		int64(uid), 
		int64(receiverid), 
		int64(Type), 
		int64(TargetType))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (m *defaultNotificationsModel) GetNotifications(ctx context.Context, uid uint64, page, pageSize int64, Type uint8) ([]*NotiItems, error) {
	var notifications []*NotiItems
	query := `
	SELECT 
		id, 
		sender_id, 
		receiver_id, 
		type, 
		content, 
		target_id, 
		is_read, 
		create_time, 
		extra, 
		target_content, 
		target_type
	FROM "public"."notifications"
	WHERE receiver_id = $1
	AND ($2 = 0 OR type = $2)
	ORDER BY create_time DESC
	LIMIT $3 OFFSET $4`
	
	offset := (page - 1) * pageSize
	err := m.QueryRowsNoCacheCtx(ctx, &notifications, query, 
		int64(uid), 
		int64(Type), 
		pageSize, 
		offset)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (m *defaultNotificationsModel) GetUnreadCount(ctx context.Context, uid uint64) (int64, error) {
	var count int64
	query := `
	SELECT COUNT(*) 
	FROM "public"."notifications"
	WHERE receiver_id = $1 
	AND is_read = false`
	
	err := m.QueryRowNoCacheCtx(ctx, &count, query, int64(uid))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *defaultNotificationsModel) MarkAllAsRead(ctx context.Context, uid uint64) error {
	query := `
	UPDATE "public"."notifications"
	SET is_read = true
	WHERE receiver_id = $1 
	AND is_read = false`
	
	_, err := m.ExecNoCacheCtx(ctx, query, int64(uid))
	return err
}