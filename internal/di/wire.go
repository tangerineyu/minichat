//go:build wireinject
// +build wireinject

package di

import (
	friendHandler "minichat/internal/handler/friend"
	friendApplyHandler "minichat/internal/handler/friend_apply"
	userHandler "minichat/internal/handler/user"
	friendRepo "minichat/internal/repo/friend"
	friendApplyRepo "minichat/internal/repo/friend_apply"
	userRepo "minichat/internal/repo/user"
	friendService "minichat/internal/service/friend"
	friendApplyService "minichat/internal/service/friend_apply"
	userService "minichat/internal/service/user"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeApp(database *gorm.DB) (*HandlerProvider, error) {
	wire.Build(
		// configs
		provideAppConfig,
		ProvideOSS,

		// repos
		userRepo.NewUserRepo,
		friendApplyRepo.NewFriendApplyRepo,
		friendRepo.NewFriendRepo,

		// services
		userService.NewUserService,
		friendApplyService.NewFriendApplyService,
		friendService.NewFriendService,

		// handlers
		userHandler.NewUserHandler,
		friendApplyHandler.NewUserHandler,
		friendHandler.NewFriendHandler,

		// provider set
		wire.Struct(new(HandlerProvider), "*"),
	)
	return &HandlerProvider{}, nil
}
