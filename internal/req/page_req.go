package req

// CursorPageReq 通用滚动分页请求参数。
// cursor: 上一页返回的 next_cursor，首次请求可为空。
// limit: 期望返回条数。
type CursorPageReq struct {
	Cursor string `form:"cursor" json:"cursor"`
	Limit  int    `form:"limit" json:"limit"`
}
