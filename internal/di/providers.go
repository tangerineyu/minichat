package di

import userHandler "minichat/internal/handler/user"

type HandlerProvider struct {
	UserHandler *userHandler.UserHandler
}
