//go:build wireinject
// +build wireinject

package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	groupHandler "minichat/internal/handler/group"
	groupApplyHandler "minichat/internal/handler/group_apply"
	messageHandler "minichat/internal/handler/message"
	userHandler "minichat/internal/handler/user"
	wsHandler "minichat/internal/handler/ws"
	friendRepo "minichat/internal/repo/friend"
	friendApplyRepo "minichat/internal/repo/friend_apply"
	groupRepo "minichat/internal/repo/group"
	groupApplyRepo "minichat/internal/repo/group_apply"
	groupMemberRepo "minichat/internal/repo/group_member"
	messageRepo "minichat/internal/repo/message"
	userRepo "minichat/internal/repo/user"
	friendService "minichat/internal/service/friend"
	friendApplyService "minichat/internal/service/friend_apply"
	groupService "minichat/internal/service/group"
	groupApplyService "minichat/internal/service/group_apply"
	messageService "minichat/internal/service/message"
	userService "minichat/internal/service/user"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var UserSet = wire.NewSet(
	// config/oss
	provideAppConfig,
	ProvideOSS,

	// repo
	userRepo.NewUserRepo,
	wire.Bind(new(userRepo.UserRepoInterface), new(*userRepo.UserRepo)),
	// service/handler
	userService.NewUserService,
	wire.Bind(new(userService.UserServiceInterface), new(*userService.UserService)),
	userHandler.NewUserHandler,
)

var FriendSet = wire.NewSet(
	// repo
	friendRepo.NewFriendRepo,
	wire.Bind(new(friendRepo.FriendRepoInterface), new(*friendRepo.FriendRepo)),

	// service/handler
	friendService.NewFriendService,
	wire.Bind(new(friendService.FriendServiceInterface), new(*friendService.FriendService)),
	friendHandler.NewFriendHandler,
)

var FriendApplySet = wire.NewSet(
	// repo
	friendApplyRepo.NewFriendApplyRepo,
	wire.Bind(new(friendApplyRepo.FriendApplyRepoInterface), new(*friendApplyRepo.FriendApplyRepo)),

	// service/handler
	friendApplyService.NewFriendApplyService,
	wire.Bind(new(friendApplyService.FriendApplyServiceInterface), new(*friendApplyService.FriendApplyService)),
	friendApplyHandler.NewUserHandler,
)

var GroupSet = wire.NewSet(
	// repo
	groupRepo.NewGroupRepo,
	wire.Bind(new(groupRepo.GroupRepoInterface), new(*groupRepo.GroupRepo)),

	// service/handler
	groupService.NewGroupService,
	//wire.Bind(new(groupService.GroupServiceInterface), new(*groupService.GroupService)),
	groupHandler.NewGroupHandler,
)
var GroupMemberSet = wire.NewSet(
	// repo
	groupMemberRepo.NewGroupMemberRepo,
	//wire.Bind(new(groupMemberRepo.GroupMemberRepoInterface), new(*groupMemberRepo.GroupMemberRepo)),
)
var GroupApplySet = wire.NewSet(
	// repo
	groupApplyRepo.NewGroupApplyRepo,

	groupApplyService.NewGroupApplyService,

	groupApplyHandler.NewGroupApplyHandler,
)

var MessageSet = wire.NewSet(
	// repo
	messageRepo.NewMessageRepo,

	// service/handler
	messageService.NewMessageService,
	messageHandler.NewMessageHandler,
)

var WSSet = wire.NewSet(
	wsHandler.NewWSHandler,
)

var HandlerProviderSet = wire.NewSet(
	UserSet,
	FriendSet,
	FriendApplySet,
	GroupSet,
	GroupMemberSet,
	GroupApplySet,
	MessageSet,
	WSSet,
	wire.Struct(new(HandlerProvider), "*"),
)

func InitializeApp(database *gorm.DB) (*HandlerProvider, error) {
	wire.Build(HandlerProviderSet)
	return &HandlerProvider{}, nil
}
