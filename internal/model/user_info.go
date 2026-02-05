package model

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	Username  string `gorm:"type:VARCHAR(64);uniqueIndex;not null" json:"username"`
	Telephone string `gorm:"type:VARCHAR(20);uniqueIndex;not null" json:"telephone"`
	Password  string `gorm:"type:VARCHAR(255);not null" json:"-"`
	Nickname  string `gorm:"type:VARCHAR(64)" json:"nickname"`
	Avatar    string `gorm:"type:VARCHAR(255)" json:"avatar"`
	Status    int8   `gorm:"type:TINYINT;not null;default:0;comment:用户状态(0离线/1在线/2注销)" json:"status"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`

	UsernameChangedAt *time.Time `gorm:"column:username_changed_at;comment:username上次修改时间" json:"username_changed_at"`
}

// TableName 定义表名。
func (User) TableName() string { return "users" }
