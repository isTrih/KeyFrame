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
		GetNotifications(ctx context.Context, uid uint64, page, pageSize int64, Type uint8) ([]*NotiItems, error)
		GetUnreadCount(ctx context.Context, uid uint64) (int64, error)
		GetTypeUnreadCount(ctx context.Context, uid uint64, notificationType uint8) (int64, error)
		GetNotificationsByType(ctx context.Context, uid uint64, page, pageSize uint64, Type uint8) ([]*NotiItems, uint64, error)
		GetGroupedNotifications(ctx context.Context, uid uint64, page, pageSize uint64, Type uint8) ([]map[string]interface{}, uint64, error)
		MarkAllAsRead(ctx context.Context, uid uint64) error
	}

	customNotificationsModel struct {
		*defaultNotificationsModel
	}

	Extra struct {
		SenderAvatar string `json:"sender_avatar"`
		SenderName   string `json:"sender_name"`
	}
	NotiItems struct {
		Id            int64     `db:"id"`             // 主键，自增 ID
		SenderId      int64     `db:"sender_id"`      // 发送人 ID，系统通知可为 NULL
		ReceiverId    int64     `db:"receiver_id"`    // 接收人 ID
		Type          int64     `db:"type"`           // 通知类型，例如 1点赞2评论3收藏4关注
		Content       string    `db:"content"`        // 通知显示内容，可用于消息列表展示
		TargetId      int64     `db:"target_id"`      // 通知对象的 ID
		IsRead        bool      `db:"is_read"`        // 是否已读标志
		CreateTime    time.Time `db:"create_time"`    // 通知创建时间
		Extra         string    `db:"extra"`          // 额外信息（JSON 格式），如帖子标题、评论摘要等
		TargetContent string    `db:"target_content"` // 目标内容
		TargetType    int64     `db:"target_type"`    // 目标类型 文章，评论，用户
	}
)

// NewNotificationsModel returns a model for the database table.
func NewNotificationsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotificationsModel {
	return &customNotificationsModel{
		defaultNotificationsModel: newNotificationsModel(conn, c, opts...),
	}
}

// ... 现有代码 ...

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
	SELECT EXISTS (
		SELECT 1 
		FROM "public"."notifications" 
		WHERE sender_id = $1 
		AND receiver_id = $2 
		AND type = $3 
		AND target_type = $4
		LIMIT 1
	)`

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

// GetNotificationsByType 获取特定类型通知及总数
func (m *defaultNotificationsModel) GetNotificationsByType(ctx context.Context, uid uint64, page, pageSize uint64, Type uint8) ([]*NotiItems, uint64, error) {
	var notifications []*NotiItems

	// 使用单一SQL查询获取结果和总数
	query := `
	WITH notification_data AS (
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
		AND type = $2
		ORDER BY create_time DESC
		LIMIT $3 OFFSET $4
	),
	notification_count AS (
		SELECT COUNT(*) as total
		FROM "public"."notifications"
		WHERE receiver_id = $1 
		AND type = $2
	)
	SELECT 
		n.*,
		(SELECT total FROM notification_count) as total_count
	FROM notification_data n`

	offset := (page - 1) * pageSize

	// 定义结构来接收查询结果
	type notificationWithCount struct {
		NotiItems
		TotalCount uint64 `db:"total_count"`
	}

	var results []notificationWithCount
	err := m.QueryRowsNoCacheCtx(ctx, &results, query,
		int64(uid),
		int64(Type),
		pageSize,
		offset)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有结果，返回空数组和0总数
	if len(results) == 0 {
		return []*NotiItems{}, 0, nil
	}

	// 转换结果并获取总数
	var totalCount uint64 = results[0].TotalCount
	notifications = make([]*NotiItems, len(results))
	for i, item := range results {
		notifications[i] = &NotiItems{
			Id:            item.Id,
			SenderId:      item.SenderId,
			ReceiverId:    item.ReceiverId,
			Type:          item.Type,
			Content:       item.Content,
			TargetId:      item.TargetId,
			IsRead:        item.IsRead,
			CreateTime:    item.CreateTime,
			Extra:         item.Extra,
			TargetContent: item.TargetContent,
			TargetType:    item.TargetType,
		}
	}

	return notifications, totalCount, nil
}

// GetGroupedNotifications 获取按发送者分组的通知（用于点赞和收藏）
func (m *defaultNotificationsModel) GetGroupedNotifications(ctx context.Context, uid uint64, page, pageSize uint64, Type uint8) ([]map[string]interface{}, uint64, error) {
	// 单一SQL查询同时获取分组数据和总数
	query := `
	WITH grouped_data AS (
		SELECT 
			json_agg(DISTINCT json_extract_path_text(extra::json, 'sender_name')) FILTER (WHERE json_extract_path_text(extra::json, 'sender_name') IS NOT NULL) AS sender_names,
			json_agg(DISTINCT json_extract_path_text(extra::json, 'sender_avatar')) FILTER (WHERE json_extract_path_text(extra::json, 'sender_avatar') IS NOT NULL) AS sender_avatars,
			type,
			COUNT(*) AS total,
			MAX(create_time) AS latest_time,
			EXTRACT(EPOCH FROM MAX(create_time))::bigint AS time_unix,
			DATE_TRUNC('day', create_time) AS group_date
		FROM "public"."notifications"
		WHERE receiver_id = $1 AND type = $2
		GROUP BY type, group_date
		ORDER BY latest_time DESC
		LIMIT $3 OFFSET $4
	),
	total_groups AS (
		SELECT COUNT(*) AS total_count FROM (
			SELECT DATE_TRUNC('day', create_time) AS day
			FROM "public"."notifications"
			WHERE receiver_id = $1 AND type = $2
			GROUP BY day
		) AS groups
	)
	SELECT 
		g.sender_names, 
		g.sender_avatars, 
		g.type, 
		g.total, 
		g.time_unix AS time,
		(SELECT total_count FROM total_groups) AS total_groups
	FROM grouped_data g
	ORDER BY g.latest_time DESC`

	offset := (page - 1) * pageSize

	// 定义临时结构接收结果
	type groupResult struct {
		SenderNames   string `db:"sender_names"`
		SenderAvatars string `db:"sender_avatars"`
		Type          int64  `db:"type"`
		Total         uint64 `db:"total"`
		Time          uint64 `db:"time"`
		TotalGroups   uint64 `db:"total_groups"`
	}

	var results []groupResult
	err := m.QueryRowsNoCacheCtx(ctx, &results, query,
		int64(uid),
		int64(Type),
		pageSize,
		offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询通知分组失败: %w", err)
	}

	if len(results) == 0 {
		return []map[string]interface{}{}, 0, nil
	}

	// 转换结果为所需的格式
	var totalGroups uint64
	groupData := make([]map[string]interface{}, 0, len(results))

	for _, r := range results {
		totalGroups = r.TotalGroups

		var names, avatars []string
		if err := json.Unmarshal([]byte(r.SenderNames), &names); err != nil {
			return nil, 0, fmt.Errorf("解析发送者名称失败: %w, 原始数据: %s", err, r.SenderNames)
		}
		if err := json.Unmarshal([]byte(r.SenderAvatars), &avatars); err != nil {
			return nil, 0, fmt.Errorf("解析发送者头像失败: %w, 原始数据: %s", err, r.SenderAvatars)
		}

		groupData = append(groupData, map[string]interface{}{
			"sender_name":   names,
			"sender_avatar": avatars,
			"type":          r.Type,
			"total":         r.Total,
			"time":          r.Time,
		})
	}

	return groupData, totalGroups, nil
}

// 获取特定类型的未读通知数
func (m *defaultNotificationsModel) GetTypeUnreadCount(ctx context.Context, uid uint64, notificationType uint8) (int64, error) {
	var count int64
	query := `
	SELECT COUNT(*) 
	FROM "public"."notifications"
	WHERE receiver_id = $1 
	AND type = $2
	AND is_read = false`

	err := m.QueryRowNoCacheCtx(ctx, &count, query, int64(uid), int64(notificationType))
	if err != nil {
		return 0, err
	}
	return count, nil
}
