package dto

// UserInfo 是对外返回的用户信息 DTO（用于 API 响应，不要放在 req 包）。
//
// 说明：
// - service 层可以返回 dto，避免依赖 handler/response 包。
// - handler 层再用 response.Success/Fail 做统一包装。
type UserInfo struct {
	UserId    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Telephone string `json:"telephone"`
}
