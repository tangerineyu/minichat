package friend

type FriendHandlerInterface interface {
	//拉黑好友
	BlackFriend()
	//取消拉黑好友
	UnBlackFriend()
	//删除好友
	DeleteFriend()
	//获取好友列表
	GetFriendList()
}
