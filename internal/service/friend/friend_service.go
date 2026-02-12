package friend

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"minichat/internal/dto"
	"minichat/internal/repo/friend"
	"minichat/internal/repo/user"
	"strings"

	"gorm.io/gorm"
)

type FriendService struct {
	friendRepo friend.FriendRepoInterface
	userRepo   user.UserRepoInterface
}

func NewFriendService(repo friend.FriendRepoInterface, userRepo user.UserRepoInterface) *FriendService {
	return &FriendService{
		friendRepo: repo,
		userRepo:   userRepo,
	}
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
	//排序名称
	sortName := remark
	fri, err := f.userRepo.GetUserById(ctx, Id)
	if err != nil {
		return errors.New("查询用户信息失败")
	}
	friendNickname := fri.Nickname
	if strings.TrimSpace(sortName) == "" {
		sortName = friendNickname
	}
	return f.friendRepo.UpdateFriendFields(ctx, Id, friendId, map[string]interface{}{
		"remark":    remark,
		"sort_name": strings.ToUpper(sortName),
	})
}

// 注意：全局排序 + 滚动分页必须以数据库的排序字段为准。
// 任何服务层二次排序都会破坏 cursor 分页的一致性（导致重复/漏数据）。

const (
	defaultFriendListLimit = 50
	maxFriendListLimit     = 200
)

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultFriendListLimit
	}
	if limit > maxFriendListLimit {
		return maxFriendListLimit
	}
	return limit
}

func (f *FriendService) GetFriendList(ctx context.Context, Id int64, cursor string, limit int) ([]*dto.FriendItem, string, error) {
	limit = normalizeLimit(limit)
	lastName, lastId := parseCursor(cursor)

	list, err := f.friendRepo.GetList(ctx, Id, 1, lastName, lastId, limit)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(list) > 0 {
		lastItem := list[len(list)-1]
		nextCursor = encoderCursor(lastItem.SortedName, lastItem.FriendId)
	}
	return list, nextCursor, nil
}

func (f *FriendService) GetBlackFriendList(ctx context.Context, Id int64, cursor string, limit int) ([]*dto.FriendItem, string, error) {
	limit = normalizeLimit(limit)
	lastName, lastId := parseCursor(cursor)

	list, err := f.friendRepo.GetList(ctx, Id, 2, lastName, lastId, limit)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(list) > 0 {
		lastItem := list[len(list)-1]
		nextCursor = encoderCursor(lastItem.SortedName, lastItem.FriendId)
	}
	return list, nextCursor, nil
}

type friendListCursor struct {
	LastSortedName string `json:"n"`
	LastID         int64  `json:"id"`
}

// 将 cursor 解码成最后一个好友的排序名称和ID，供下一页查询使用  解密
func parseCursor(cursor string) (lastSortedName string, lastID int64) {
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		return "", 0
	}

	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		// 容错：非法 cursor 当作第一页处理
		return "", 0
	}
	var c friendListCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return "", 0
	}
	return c.LastSortedName, c.LastID
}

// 将最后一个好友的排序名称和ID编码成 cursor，供下一页查询使用  加密
func encoderCursor(lastSortedName string, lastID int64) string {
	c := friendListCursor{LastSortedName: lastSortedName, LastID: lastID}
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
