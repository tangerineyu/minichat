package model

import "time"

// Message 私聊消息
// MsgType: 1文字, 2图片, 3语音, 4视频
type Message struct {
	ID         int64     `gorm:"primaryKey;type:BIGINT" json:"id"`
	SenderID   int64     `gorm:"type:BIGINT;not null;index:idx_sender_receiver,priority:1" json:"sender_id"`
	ReceiverID int64     `gorm:"type:BIGINT;not null;index:idx_sender_receiver,priority:2" json:"receiver_id"`
	MsgType    int8      `gorm:"type:TINYINT;not null" json:"msg_type"`
	Content    string    `gorm:"type:TEXT" json:"content"`
	IsRead     bool      `gorm:"type:BOOLEAN;not null;default:false" json:"is_read"`
	CreatedAt  time.Time `gorm:"type:DATETIME;not null;autoCreateTime;index" json:"created_at"`
}

func (Message) TableName() string { return "messages" }
