package model

// GroupApply 群申请表
// InviteId 0表示用户主动申请加入，非0表示被邀请加入
type GroupApply struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID   int64  `json:"group_id" gorm:"type:BIGINT;not null;index"`
	UserID    int64  `json:"user_id" gorm:"type:BIGINT;not null;index"`
	InviterId int64  `json:"inviter_id" gorm:"type:BIGINT;default:0;not null;index"`
	Reason    string `json:"reason" gorm:"type:VARCHAR(255)"`
	Status    int8   `json:"status" gorm:"type:TINYINT;not null;default:0;comment:0待处理/1已同意/2已拒绝"`
	HandlerId int64  `json:"handler_id" gorm:"type:BIGINT;default:0;not null;index"`

	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (GroupApply) TableName() string {
	return "group_applies"
}
