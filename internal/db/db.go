package db

import (
	"minichat/internal/model"

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
