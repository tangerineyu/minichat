package friend

import (
	"minichat/internal/repo/friend"
	"minichat/internal/req"
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

func (f *FriendService) GetFriendList(Id int64) ([]*req.FriendListReq, error) {
	//TODO implement me
	panic("implement me")
}

func NewFriendService(repo friend.FriendRepoInterface) FriendServiceInterface {
	return &FriendService{friendRepo: repo}
}
