package friend_apply

import (
	"context"
	"minichat/internal/model"
	"minichat/internal/req"
)

type FriendApplyRepoInterface interface {
	// 发送好友申请
	CreateFriendApply(ctx context.Context, formUserId, toUsrId int64, applyMsg string) (*model.FriendApply, error)
	// 通过ID查询好友申请
	GetFriendApplyById(ctx context.Context, applyId int64) (*model.FriendApply, error)
	// 更新好友申请状态
	UpdateApplyStatus(ctx context.Context, applyId int64, status int) error
	// 获取用户的好友申请列表
	GetFriendApplyList(ctx context.Context, id int64) ([]*req.FriendApplyListReq, error)
}
