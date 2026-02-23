package model

import "time"

// Message 私聊消息
// MsgType: 1文字, 2图片, 3语音, 4视频
type Message struct {
	ID         int64 `gorm:"primaryKey;type:BIGINT" json:"id"` //snowflake生成的全局唯一ID
	SenderID   int64 `gorm:"type:BIGINT;not null;index" json:"sender_id"`
	ReceiverID int64 `gorm:"type:BIGINT;not null;index" json:"receiver_id"`

	SessionType int8 `gorm:"type:TINYINT;not null" json:"session_type"` // 1私聊, 2群聊
	MsgType     int8 `gorm:"type:TINYINT;not null" json:"msg_type"`     //  决定内容格式: 1-文本, 2-图片, 3-文件, 4-系统通知文本

	Content string `gorm:"type:TEXT" json:"content"`
	Seq     int64  `gorm:"type:BIGINT;not null;index" json:"seq"` //连续递增的消息序号，用于补发消息

	Status    int8      `gorm:"type:TINYINT;not null;default:0" json:"status"` // 0未读, 1撤回, 2删除 (仅对接收方有效)
	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime;index" json:"created_at"`
}

func (Message) TableName() string { return "messages" }
