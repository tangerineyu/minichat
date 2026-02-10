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

func (f *FriendHandler) BlackFriend(c *gin.Context) {
	var in req.BlackFriendReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, 400, "")
	}

}

func (f *FriendHandler) UnBlackFriend(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (f *FriendHandler) DeleteFriend(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (f *FriendHandler) GetFriendList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

var _ (FriendHandlerInterface) = (*FriendHandler)(nil)

func NewFriendHandler(svc friend.FriendServiceInterface) *FriendHandler {
	return &FriendHandler{
		friendService: svc,
	}
}
