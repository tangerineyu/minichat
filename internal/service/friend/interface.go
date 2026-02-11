package friend

import (
	"context"
	"minichat/internal/dto"
)

type FriendServiceInterface interface {
	// 拉黑好友
	BlackFriend(ctx context.Context, Id, friendId int64) error
	// 取消拉黑好友
	UnBlackFriend(ctx context.Context, Id, friendId int64) error
	// 删除好友
	DeleteFriend(ctx context.Context, Id, friendId int64) error

	// 获取好友列表（滚动分页）
	GetFriendList(ctx context.Context, Id int64, cursor string, limit int) ([]*dto.FriendItem, string, error)
	// 修改好友备注
	UpdateFriendRemark(ctx context.Context, Id, friendId int64, remark string) error
	// 获取拉黑好友列表（滚动分页）
	GetBlackFriendList(ctx context.Context, Id int64, cursor string, limit int) ([]*dto.FriendItem, string, error)
}
