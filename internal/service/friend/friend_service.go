package friend

import (
	"context"
	"errors"
	"minichat/internal/dto"
	"minichat/internal/repo/friend"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type FriendService struct {
	friendRepo friend.FriendRepoInterface
}

func (f *FriendService) BlackFriend(ctx context.Context, Id, friendId int64) error {
	err := f.friendRepo.UpdateFriendFields(ctx, Id, friendId, map[string]interface{}{
		"status": 2,
	})
	// TODO: 可以加入redis黑名单
	return err
}

func (f *FriendService) UnBlackFriend(ctx context.Context, Id, friendId int64) error {
	return f.friendRepo.UpdateFriendFields(ctx, Id, friendId, map[string]interface{}{
		"status": 1,
	})
}

func (f *FriendService) DeleteFriend(ctx context.Context, Id, friendId int64) error {
	//TODO implement me
	panic("implement me")
}

// 修改好友备注
func (f *FriendService) UpdateFriendRemark(ctx context.Context, Id, friendId int64, remark string) error {
	relation, err := f.friendRepo.GetFriendRelation(ctx, Id, friendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("你们还不是好友关系，请先添加好友")
		}
		return errors.New("查询好友关系失败")
	}
	if relation.Status != 1 {
		return errors.New("你们还不是好友关系，请先添加好友")
	}
	if len(remark) > 64 {
		return errors.New("备注长度不能超过64个字符")
	}
	return f.friendRepo.UpdateFriendFields(ctx, Id, friendId, map[string]interface{}{
		"remark": remark,
	})
}

func displayNameForSort(it *dto.FriendItem) string {
	if it == nil {
		return ""
	}
	name := strings.TrimSpace(it.FriendRemark)
	if name == "" {
		name = strings.TrimSpace(it.FriendNickname)
	}
	return name
}

func sortFriendItemsByRemarkThenNickname(list []*dto.FriendItem) {
	sort.SliceStable(list, func(i, j int) bool {
		a := strings.ToUpper(displayNameForSort(list[i]))
		b := strings.ToUpper(displayNameForSort(list[j]))
		if a == b {
			// 次级排序：friend_id，保证稳定可预期
			if list[i] == nil || list[j] == nil {
				return a < b
			}
			return list[i].FriendId < list[j].FriendId
		}
		return a < b
	})
}

func (f *FriendService) GetFriendList(ctx context.Context, Id int64) ([]*dto.FriendItem, error) {
	list, err := f.friendRepo.GetList(ctx, Id, 1)
	if err != nil {
		return nil, err
	}
	sortFriendItemsByRemarkThenNickname(list)
	return list, nil
}

func (f *FriendService) GetBlackFriendList(ctx context.Context, Id int64) ([]*dto.FriendItem, error) {
	list, err := f.friendRepo.GetList(ctx, Id, 2)
	if err != nil {
		return nil, err
	}
	sortFriendItemsByRemarkThenNickname(list)
	return list, nil
}

func NewFriendService(repo friend.FriendRepoInterface) FriendServiceInterface {
	return &FriendService{friendRepo: repo}
}
