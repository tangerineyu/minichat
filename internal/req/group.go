package req

type CreateGroupReq struct {
	Name         string `json:"name" validate:"required"`
	Announcement string `json:"announcement"`
	Avatar       string `json:"avatar"`
}

type UpdateGroupInfoReq struct {
	Name         string `json:"name"`
	Announcement string `json:"announcement"`
	Avatar       string `json:"avatar"`
}

type AddGroupMembersReq struct {
	UserIds []int64 `json:"user_ids"`
}

type DealWithGroupApplyReq struct {
	ApplyId int64 `json:"apply_id"`
	Status  int8  `json:"status" comment:"待处理0, 同意 1，拒绝 2"`
}
