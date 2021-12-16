package initialize

import (
	"gin-derived/global"
	DbClient "gin-derived/pkg/gorm"
	"github.com/jinzhu/gorm"
)

func InitDB() *gorm.DB {
	options := &global.GCONFIG.Database
	g := DbClient.GetGorm("default", options)
	return g.Db
}
