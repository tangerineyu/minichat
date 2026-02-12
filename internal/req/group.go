package req

type CreateGroupReq struct {
	Name         string `json:"name" validate:"required"`
	Announcement string `json:"announcement"`
	Avatar       string `json:"avatar"`
}
