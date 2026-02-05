package model

import "time"

type Group struct {
	ID           int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	Name         string `gorm:"type:VARCHAR(64);not null" json:"name"`
	OwnerID      int64  `gorm:"type:BIGINT;not null;index" json:"owner_id"`
	Announcement string `gorm:"type:TEXT" json:"announcement"`
	Avatar       string `gorm:"type:VARCHAR(255)" json:"avatar"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (Group) TableName() string { return "groups" }

// GroupMember 群成员
// Role: 0普通成员, 1管理员, 2群主
type GroupMember struct {
	GroupID  int64  `gorm:"type:BIGINT;not null;primaryKey" json:"group_id"`
	UserID   int64  `gorm:"type:BIGINT;not null;primaryKey" json:"user_id"`
	Nickname string `gorm:"type:VARCHAR(64)" json:"nickname"`
	Role     int8   `gorm:"type:TINYINT;not null;default:0" json:"role"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (GroupMember) TableName() string { return "group_members" }
