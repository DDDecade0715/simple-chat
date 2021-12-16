package main

import (
	"gin-derived/api/models"
	"gin-derived/global"
	"gin-derived/initialize"
	"log"
	"os"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	global.GVIPER = initialize.InitViper(dir + "/../config.yaml") //初始化viper
	global.GDB = initialize.InitDB()                              //初始化数据库
}
func main() {
	// 执行自动迁移
	global.GDB.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.ChatImage{},
		&models.UserRelation{},
		&models.Group{},
		&models.GroupMember{},
		&models.SystemMessage{},
	)
}
