//go:build wireinject
// +build wireinject

package di

import (
	friendApplyHandler "minichat/internal/handler/friend_apply"
	userHandler "minichat/internal/handler/user"
	friendApplyRepo "minichat/internal/repo/friend_apply"
	userRepo "minichat/internal/repo/user"
	friendApplyService "minichat/internal/service/friend_apply"
	userService "minichat/internal/service/user"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeApp(database *gorm.DB) (*HandlerProvider, error) {
	wire.Build(
		// repos
		userRepo.NewUserRepo,
		friendApplyRepo.NewFriendApplyRepo,

		// services
		userService.NewUserService,
		friendApplyService.NewFriendApplyService,

		// handlers
		userHandler.NewUserHandler,
		friendApplyHandler.NewUserHandler,

		// provider set
		wire.Struct(new(HandlerProvider), "*"),
	)
	return &HandlerProvider{}, nil
}
