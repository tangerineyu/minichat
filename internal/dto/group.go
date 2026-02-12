package dto

type GroupInfo struct {
	GroupId      int64  `json:"group_id"`
	GroupName    string `json:"group_name"`
	Announcement string `json:"announcement"`
	Avatar       string `json:"avatar"`
}
