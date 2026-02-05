package main

import (
	"minichat/internal/config"
	"minichat/internal/db"
	"minichat/internal/di"
	userHandler "minichat/internal/handler/user"
	"minichat/internal/repo/user"
	"minichat/internal/router"
	userService "minichat/internal/service/user"
	"minichat/util/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.GetConfig()
	logger.InitLogger(&logger.Config{
		Path:       cfg.Logger.Path,
		Level:      cfg.Logger.Level,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		ToStdout:   cfg.Logger.ToStdout,
	})
	defer logger.Sync()

	gin.SetMode(cfg.Server.GinMode)

	// InitDB reads cfg.DB.driver/cfg.DB.dsn and runs AutoMigrate.
	database, err := db.InitDB()
	if err != nil {
		zap.L().Fatal("open db failed", zap.Error(err))
	}

	repo := user.NewUserRepo(database)
	svc := userService.NewUserService(repo)
	h := userHandler.NewUserHandler(svc)

	providers := &di.HandlerProvider{UserHandler: h}

	r := gin.New()
	router.SetupRouter(r, providers)

	zap.L().Info("http server starting", zap.String("addr", cfg.Server.Addr), zap.String("mode", cfg.Server.GinMode))
	if err := r.Run(cfg.Server.Addr); err != nil {
		zap.L().Fatal("http server stopped", zap.Error(err))
	}
}
