package friend

import (
	"minichat/internal/service/friend"

	"github.com/gin-gonic/gin"
)

type FriendHandler struct {
	friendService friend.FriendServiceInterface
}

func (f *FriendHandler) BlackFriend(c *gin.Context) {
	//TODO implement me
	panic("implement me")

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
