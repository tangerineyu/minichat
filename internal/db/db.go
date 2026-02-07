package db

import (
	"minichat/internal/config"
	"minichat/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenForDev 打开一个本地 sqlite 数据库，便于快速跑通注册/登录。
// 你接入 MySQL 时可以把这个替换成 mysql driver + DSN。
func OpenForDev() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("minichat.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	return db, nil
}
func InitDB() (*gorm.DB, error) {
	cfg := config.GetConfig()
	var dialector gorm.Dialector
	switch cfg.DB.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DB.DSN)
	case "mysql":
		dialector = mysql.Open(cfg.DB.DSN)
	default:
		dialector = mysql.Open("minichat.db")
	}
	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err == nil {
		// 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(100)
		// 设置连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.FriendApply{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
