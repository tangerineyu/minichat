package event

const (
	EventMessageNew      EventType = "chat.message.new"      // 新消息事件
	EventMessageWithdraw EventType = "chat.message.withdraw" // 消息撤回事件
	EventChatRead        EventType = "chat.message.read"     // 消息已读事件
	EventChatTyping      EventType = "chat.message.typing"   // 消息输入中事件

	EventGroupApplyed   EventType = "group.apply"     // 群申请事件
	EventGroupApproved  EventType = "group.approve"   // 群批准事件
	EventGroupRejected  EventType = "group.rejected"  // 群拒绝事件
	EventGroupDissolved EventType = "group.dissolved" // 群解散事件
	EventGroupKicked    EventType = "group.kick"      // 群踢人事件

	EventSysKickout EventType = "sys.kickout" // 多端登录被踢事件

	EventFriendApplyNew EventType = "friend.apply.new" // 新好友申请事件
)
