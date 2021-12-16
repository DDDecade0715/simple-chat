package gorm

import (
	"fmt"
	"gin-derived/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Gorm struct {
	Db     *gorm.DB
	Driver string
	Dsn    string
}

var GormPool = make(map[string]*Gorm)

func GetGorm(name string, gorm *config.Database) *Gorm {
	if p, ok := GormPool[name]; ok {
		return p
	}
	p, err := newGorm(gorm)
	if err != nil {
		fmt.Printf("new Gorm failed: %s \n", err)
	}
	GormPool[name] = p
	return p
}

func newGorm(cfg *config.Database) (*Gorm, error) {
	var dsn string
	dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Protocol,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	db, err := gorm.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	db.DB().SetConnMaxLifetime(cfg.MaxLifetime) //最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(cfg.MaxOpens)       //设置最大连接数
	db.DB().SetMaxIdleConns(cfg.MaxIdles)       //设置闲置连接数
	db.SingularTable(true)                      //设置全局表名禁用复数
	if cfg.RunMode == "debug" {
		db.LogMode(true)
	}
	p := &Gorm{
		Db:     db,
		Driver: cfg.Driver,
		Dsn:    dsn,
	}
	return p, nil
}
