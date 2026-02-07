package model

import "time"

type FriendApply struct {
	ID int64 `gorm:"primaryKey;type:BIGINT" json:"id"`

	FromUserId string `gorm:"not null;index:idx_from_user,priority:1" json:"from_user_id"`
	ToUserId   string `gorm:"not null;index:idx_to_user,priority:2" json:"to_user_id"`
	ApplyMsg   string `gorm:"type:VARCHAR(255)" json:"apply_msg"`
	Status     int8   `gorm:"type:TINYINT;not null;default:0;comment:0申请中/1已通过/2拒绝" json:"status"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (FriendApply) TableName() string { return "friend_apply" }
