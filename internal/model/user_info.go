package model

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey;type:BIGINT" json:"id"`
	UserId    string `gorm:"type:VARCHAR(64);uniqueIndex;not null" json:"user_id"`
	Telephone string `gorm:"type:VARCHAR(20);uniqueIndex;not null" json:"telephone"`
	Password  string `gorm:"type:VARCHAR(255);not null" json:"-"`
	Nickname  string `gorm:"type:VARCHAR(64)" json:"nickname"`
	Avatar    string `gorm:"type:VARCHAR(255)" json:"avatar"`
	Status    int8   `gorm:"type:TINYINT;not null;default:0;comment:用户状态(0离线/1在线/2注销)" json:"status"`

	CreatedAt time.Time `gorm:"type:DATETIME;not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME;not null;autoUpdateTime" json:"updated_at"`

	UserIdChangedAt *time.Time `gorm:"column:user_id_changed_at;comment:userid上次修改时间" json:"user_id_changed_at"`
}

// TableName 定义表名。
func (User) TableName() string { return "users" }
