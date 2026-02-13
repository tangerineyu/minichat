package dto

type GroupInfo struct {
	GroupId      int64  `json:"group_id"`
	GroupName    string `json:"group_name"`
	Announcement string `json:"announcement"`
	Avatar       string `json:"avatar"`
}

type GroupDetailInfo struct {
	GroupInfo
	Members []*GroupMemberInfo `json:"members"`
}

// GroupAddStatusCode 表示邀请加群的处理结果（注意：这不是 http code）
// 0: 直接加入成功（管理员/群主邀请）
// 1: 已提交申请，等待管理员/群主审核（普通成员邀请）
type GroupAddStatusCode int8

const (
	GroupAddJoined          GroupAddStatusCode = 0
	GroupAddPendingApproval GroupAddStatusCode = 1
)

type GroupAddStatus struct {
	Status       GroupAddStatusCode `json:"status"`
	AddedUserIds []int64            `json:"added_user_ids,omitempty"`
	ApplyIds     []int64            `json:"apply_ids,omitempty"`
	Msg          string             `json:"msg,omitempty"`
}
