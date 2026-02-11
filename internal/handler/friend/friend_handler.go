package friend

import (
	"minichat/internal/handler/response"
	"minichat/internal/req"
	"minichat/internal/service/friend"

	"github.com/gin-gonic/gin"
)

type FriendHandler struct {
	friendService friend.FriendServiceInterface
}

func (f *FriendHandler) UpdateFriendRemark(c *gin.Context) {
	var in req.UpdateFriendRemarkReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	id := c.GetInt64("id")
	if err := f.friendService.UpdateFriendRemark(c.Request.Context(), id, in.FriendId, in.Remark); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

func (f *FriendHandler) GetBlackFriendList(c *gin.Context) {
	id := c.GetInt64("id")
	var page req.CursorPageReq
	_ = c.ShouldBindQuery(&page)
	list, nextCursor, err := f.friendService.GetBlackFriendList(c.Request.Context(), id, page.Cursor, page.Limit)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list, "next_cursor": nextCursor})
}

func (f *FriendHandler) BlackFriend(c *gin.Context) {
	var in req.BlackFriendReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	id := c.GetInt64("id")
	if id == in.FriendId {
		response.Fail(c, 400, "不能拉黑自己")
		return
	}
	if err := f.friendService.BlackFriend(c.Request.Context(), id, in.FriendId); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, nil)
}

func (f *FriendHandler) UnBlackFriend(c *gin.Context) {
	var in req.BlackFriendReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	id := c.GetInt64("id")
	if err := f.friendService.UnBlackFriend(c.Request.Context(), id, in.FriendId); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, nil)
}

func (f *FriendHandler) DeleteFriend(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (f *FriendHandler) GetFriendList(c *gin.Context) {
	id := c.GetInt64("id")
	var page req.CursorPageReq
	_ = c.ShouldBindQuery(&page)
	list, nextCursor, err := f.friendService.GetFriendList(c.Request.Context(), id, page.Cursor, page.Limit)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list, "next_cursor": nextCursor})
}

var _ (FriendHandlerInterface) = (*FriendHandler)(nil)

func NewFriendHandler(svc friend.FriendServiceInterface) *FriendHandler {
	return &FriendHandler{
		friendService: svc,
	}
}
