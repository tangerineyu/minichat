package friend

import "context"

type FriendRepoInterface interface {
	// 成为好友
	MakeFriends(ctx context.Context, applyId int64) error
}
