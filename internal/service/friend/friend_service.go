package friend

import "minichat/internal/repo/friend"

type FriendService struct {
	friendRepo friend.FriendRepoInterface
}

func NewFriendService(repo friend.FriendRepoInterface) FriendServiceInterface {
	return &FriendService{friendRepo: repo}
}
