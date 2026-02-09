package friend_apply

import (
	"context"
	"minichat/internal/model"
	"minichat/internal/req"
	"minichat/util/snowflake"

	"gorm.io/gorm"
)

var _ FriendApplyRepoInterface = (*FriendApplyRepo)(nil)

type FriendApplyRepo struct {
	db *gorm.DB
}

func (f *FriendApplyRepo) GetFriendApplyList(ctx context.Context, id int64) ([]*req.FriendApplyListReq, error) {
	var list []*req.FriendApplyListReq
	err := f.db.WithContext(ctx).Table("friend_apply").
		Where("friend_apply.id as apply_id, "+
			"friend_apply.apply_msg as apply_msg, "+
			"friend_apply.status as status, "+
			"friend_apply.from_user_id as from_user_id, "+
			"user.nickname as from_user_nickname"+
			"user.avatar as from_user_avatar"+
			"friend_apply.created_at as created_at").
		Joins("left join users on friend_apply.from_user_id = user.id").
		Where("friend_apply.apply_id = ?", id).
		Order("created_at desc").
		Scan(&list).Error
	return list, err

}

func (f *FriendApplyRepo) UpdateApplyStatus(ctx context.Context, applyId int64, status int) error {
	return f.db.WithContext(ctx).Model(model.FriendApply{}).
		Where("id = ?", applyId).
		Update("status", status).Error
}

func (f *FriendApplyRepo) GetFriendApplyById(ctx context.Context, applyId int64) (*model.FriendApply, error) {
	var apply model.FriendApply
	if err := f.db.WithContext(ctx).Where("id = ?", applyId).First(&apply).Error; err != nil {
		return nil, err
	}
	return &apply, nil
}

func (f *FriendApplyRepo) CreateFriendApply(ctx context.Context, formUserId, toUsrId int64, applyMsg string) (*model.FriendApply, error) {
	Id := snowflake.GenInt64ID()
	apply := &model.FriendApply{
		ID:         Id,
		FromUserId: formUserId,
		ToUserId:   toUsrId,
		ApplyMsg:   applyMsg,
		Status:     0,
	}
	if err := f.db.WithContext(ctx).Create(apply).Error; err != nil {
		return nil, err
	}
	return apply, nil
}

func NewFriendApplyRepo(db *gorm.DB) FriendApplyRepoInterface {
	return &FriendApplyRepo{db: db}
}
