package friend

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/model"
)

type FriendRepoInterface interface {
	// 成为好友
	MakeFriends(ctx context.Context, applyId int64, a2rNickname, r2aNickname string) error
	// 删除好友
	DeleteFriend(ctx context.Context, Id, friendId int64) error
	// 获取好友/黑名单列表
	GetList(ctx context.Context, Id int64, status int8, lastSortedName string, lastId int64, limit int) ([]*dto.FriendItem, error)
	// 拉黑，解除拉黑，修改好友备注等操作
	UpdateFriendFields(ctx context.Context, Id, friendId int64, fields map[string]interface{}) error
	// 获取好友关系状态，0-非好友，1-好友，2-黑名单
	GetFriendRelation(ctx context.Context, Id, friendId int64) (*model.Friend, error)
}
