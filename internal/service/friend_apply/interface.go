package friend_apply

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/req"
)

type FriendApplyServiceInterface interface {
	// 发送好友申请
	SendFriendApply(ctx context.Context, fromUserId int64, in req.SendFriendApplyReq) error
	// 处理好友申请，接受或拒绝
	DealWithFriendApply(ctx context.Context, Id int64, in req.DealWithFriendApplyReq) error
	// 获取好友申请列表
	GetFriendApply(ctx context.Context, id int64) ([]*dto.FriendApplyItem, error)
}
