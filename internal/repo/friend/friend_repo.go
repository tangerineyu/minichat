package friend

import (
	"context"
	"minichat/internal/dto"
	"minichat/internal/model"

	"gorm.io/gorm"
)

type FriendRepo struct {
	db *gorm.DB
}

func (f *FriendRepo) GetFriendRelation(ctx context.Context, Id, friendId int64) (*model.Friend, error) {
	var relation *model.Friend
	err := f.db.WithContext(ctx).
		Where("user_id = ? AND friend_id = ?", Id, friendId).
		First(&relation).Error
	if err != nil {
		return nil, err
	}
	return relation, nil
}

func (f *FriendRepo) UpdateFriendFields(ctx context.Context, Id, friendId int64, fields map[string]interface{}) error {
	return f.db.WithContext(ctx).Model(model.Friend{}).
		Where("user_id = ? AND friend_id = ?", Id, friendId).
		Updates(fields).Error
}

func (f *FriendRepo) DeleteFriend(ctx context.Context, Id, friendId int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *FriendRepo) GetList(ctx context.Context, Id int64, status int8) ([]*dto.FriendItem, error) {
	var list []*dto.FriendItem
	err := f.db.WithContext(ctx).Table("friends").
		Select("friends.friend_id as friend_id, "+
			"friends.remark as friend_remark, "+
			"users.nickname as friend_nickname, "+
			"users.avatar as friend_avatar, "+
			"friends.created_at as create_at").
		Joins("left join users on friends.friend_id = users.id").
		Where("friends.user_id = ? AND friends.status = ?", Id, status).
		Order("friends.created_at desc").
		Scan(&list).Error
	return list, err
}

func (f *FriendRepo) MakeFriends(ctx context.Context, applyId int64, a2rRemark, r2aRemark string) error {
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
			Remark:   r2aRemark,
			Status:   1,
		}
		friend2 := &model.Friend{
			UserId:   apply.ToUserId,
			FriendId: apply.FromUserId,
			Remark:   a2rRemark,
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
