package main

import (
	"minichat/internal/db"
	userHandler "minichat/internal/handler/user"
	"minichat/internal/repo/user"
	"minichat/internal/router"
	userService "minichat/internal/service/user"
	"minichat/util/logger"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger(nil)
	defer logger.Sync()

	database, err := db.OpenForDev()
	if err != nil {
		zap.L().Fatal("open db failed", zap.Error(err))
	}

	repo := user.NewUserRepo(database)
	svc := userService.NewUserService(repo)
	h := userHandler.NewUserHandler(svc)

	r := router.NewRouter(h)

	addr := os.Getenv("MINICHAT_LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	zap.L().Info("http server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		zap.L().Fatal("http server stopped", zap.Error(err))
	}
}
