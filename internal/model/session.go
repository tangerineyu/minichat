package model

import "time"

// Session 会话列表（Inbox）
// Type: 1私聊, 2群聊
type Session struct {
	ID          int64     `gorm:"primaryKey;type:BIGINT" json:"id"`
	UserID      int64     `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_target_type,priority:1" json:"user_id"`
	TargetID    int64     `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_target_type,priority:2" json:"target_id"`
	Type        int8      `gorm:"type:TINYINT;not null;uniqueIndex:uk_user_target_type,priority:3" json:"type"`
	LastMsg     string    `gorm:"type:VARCHAR(255)" json:"last_msg"`
	UnreadCount int32     `gorm:"type:INT;not null;default:0" json:"unread_count"`
	UpdatedAt   time.Time `gorm:"type:DATETIME;not null;autoUpdateTime;index" json:"updated_at"`
}

func (Session) TableName() string { return "sessions" }
