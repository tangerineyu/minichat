package friend

import (
	"errors"
	"minichat/internal/dto"
	"minichat/internal/repo/friend"

	"go.uber.org/zap"
)

type FriendService struct {
	friendRepo friend.FriendRepoInterface
}

func (f *FriendService) BlackFriend(Id, friendId int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *FriendService) UnBlackFriend(Id, friendId int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *FriendService) DeleteFriend(Id, friendId int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *FriendService) GetFriendList(Id int64) ([]*dto.FriendItem, error) {
	zap.L().Warn("GetFriendList not implemented")
	return nil, errors.New("获取好友列表暂未实现")
}

func NewFriendService(repo friend.FriendRepoInterface) FriendServiceInterface {
	return &FriendService{friendRepo: repo}
}
