package friend

import (
	"context"
	"minichat/internal/model"

	"gorm.io/gorm"
)

type FriendRepo struct {
	db *gorm.DB
}

func (f *FriendRepo) MakeFriends(ctx context.Context, applyId int64) error {
	// 这里可以使用事务来确保两个用户的好友关系同时创建成功
	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var apply model.FriendApply
		// 根据好友申请ID查询好友申请记录
		if err := tx.First(&apply, applyId).Error; err != nil {
			return err
		}
		// 更新好友申请状态为已接受
		if err := tx.Model(&apply).Update("status", 1).Error; err != nil {
			return err
		}
		friend1 := &model.Friend{
			UserId:   apply.FromUserId,
			FriendId: apply.ToUserId,
			// TODO: 可以设置备注，默认是好友的昵称
			Status: 1,
		}
		friend2 := &model.Friend{
			UserId:   apply.ToUserId,
			FriendId: apply.FromUserId,
			Status:   1,
		}
		friends := []model.Friend{*friend1, *friend2}
		if err := tx.Create(&friends).Error; err != nil {
			// 如果因为唯一索引冲突（已经是好友），这里会报错并导致回滚
			return err
		}
		return nil
	})
}

func NewFriendRepo(db *gorm.DB) FriendRepoInterface {
	return &FriendRepo{db: db}
}
