//go:build wireinject
// +build wireinject

package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	groupHandler "minichat/internal/handler/group"
	userHandler "minichat/internal/handler/user"
	friendRepo "minichat/internal/repo/friend"
	friendApplyRepo "minichat/internal/repo/friend_apply"
	groupRepo "minichat/internal/repo/group"
	userRepo "minichat/internal/repo/user"
	friendService "minichat/internal/service/friend"
	friendApplyService "minichat/internal/service/friend_apply"
	groupService "minichat/internal/service/group"
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
	wire.Bind(new(groupService.GroupServiceInterface), new(*groupService.GroupService)),
	groupHandler.NewGroupHandler,
)

var HandlerProviderSet = wire.NewSet(
	UserSet,
	FriendSet,
	FriendApplySet,
	GroupSet,
	wire.Struct(new(HandlerProvider), "*"),
)

func InitializeApp(database *gorm.DB) (*HandlerProvider, error) {
	wire.Build(HandlerProviderSet)
	return &HandlerProvider{}, nil
}
