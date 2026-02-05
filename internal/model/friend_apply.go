package model

import "time"

// Friend 表示好友关系（推荐插入两条：A->B 与 B->A）
// Status: 0申请中, 1已通过, 2黑名单
type Friend struct {
	ID       int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	UserID   int64  `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_friend" json:"user_id"`
	FriendID int64  `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_friend" json:"friend_id"`
	Remark   string `gorm:"type:VARCHAR(64)" json:"remark"`
	Status   int8   `gorm:"type:TINYINT;not null;default:0;comment:0申请中/1已通过/2黑名单" json:"status"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (Friend) TableName() string { return "friends" }
