package friend

import "minichat/internal/req"

type FriendServiceInterface interface {
	// 拉黑好友
	BlackFriend(Id, friendId int64) error
	// 取消拉黑好友
	UnBlackFriend(Id, friendId int64) error
	// 删除好友
	DeleteFriend(Id, friendId int64) error
	// 获取好友列表
	GetFriendList(Id int64) ([]*req.FriendListReq, error)
}
