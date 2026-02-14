package model

import "time"

type Group struct {
	ID   int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	Name string `gorm:"type:VARCHAR(64);not null" json:"name"`
	// 群主ID，关联用户表的ID
	OwnerID      int64  `gorm:"type:BIGINT;not null;index" json:"owner_id"`
	Announcement string `gorm:"type:TEXT" json:"announcement"`
	Avatar       string `gorm:"type:VARCHAR(255)" json:"avatar"`
	// JoinMode 入群模式：0-直接加入, 1-需审核, 2-禁止加入
	JoinMode  int8      `gorm:"type:TINYINT;not null;default:0" json:"join_mode"`
	Status    int8      `gorm:"type:TINYINT;not null;default:0;comment:0正常/1解散" json:"status"`
	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (Group) TableName() string { return "groups" }

// GroupMember 群成员
// Role: 0普通成员, 1管理员, 2群主
// IsMuted: 是否被禁言
type GroupMember struct {
	GroupID  int64  `gorm:"type:BIGINT;not null;primaryKey" json:"group_id"`
	UserID   int64  `gorm:"type:BIGINT;not null;primaryKey" json:"user_id"`
	Nickname string `gorm:"type:VARCHAR(64)" json:"nickname"`
	Role     int8   `gorm:"type:TINYINT;not null;default:0" json:"role"`

	IsMuted   bool       `gorm:"type:TINYINT;not null;default:false" json:"is_muted"`
	MutedTime *time.Time `gorm:"type:DATETIME" json:"muted_time"`

	Status int8 `gorm:"type:TINYINT;not null;default:0;comment:0正常/1禁言/2已经退出" json:"status"`
	// InviterID 0表示用户主动加入，非0表示被邀请加入
	InviterID int64 `gorm:"type:BIGINT;default:0" json:"inviter_id"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`
}

func (GroupMember) TableName() string { return "group_members" }
