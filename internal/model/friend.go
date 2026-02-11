package model

import "time"

// Friend 表示好友关系（推荐插入两条：A->B 与 B->A）
// Status: 1好友, 2黑名单
type Friend struct {
	Id       int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	UserId   int64  `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_friend" json:"user_id"`
	FriendId int64  `gorm:"type:BIGINT;not null;index;uniqueIndex:uk_user_friend" json:"friend_id"`
	Remark   string `gorm:"type:VARCHAR(64)" json:"remark" comment:"好友备注"`
	Status   int8   `gorm:"type:TINYINT;not null;default:0;comment:1好友/2黑名单" json:"status"`
	// 用于排序
	SortName  string    `gorm:"type:VARCHAR(255);not null;default:''" json:"sort_name"`
	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (Friend) TableName() string { return "friends" }
