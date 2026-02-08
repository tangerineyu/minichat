package friend_apply

import (
	"context"
	"minichat/internal/model"
)

type FriendApplyRepoInterface interface {
	// 发送好友申请
	CreateFriendApply(ctx context.Context, formUserId, toUsrId int64, applyMsg string) (*model.FriendApply, error)
	// 通过ID查询好友申请
	GetFriendApplyById(ctx context.Context, applyId int64) (*model.FriendApply, error)
	// 更新好友申请状态
	UpdateApplyStatus(ctx context.Context, applyId int64, status int) error
}
