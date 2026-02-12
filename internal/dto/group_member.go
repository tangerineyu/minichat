package dto

type GroupMemberInfo struct {
	GroupId  int64  `json:"group_id"`
	MemberId int64  `json:"member_id"`
	Nickname string `json:"nickname"`
	Role     int8   `json:"role"`
	Status   int8   `json:"status"`
}
