package initialize

import (
	"fmt"
	"gin-derived/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
)

func InitViper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		config = "./config.yaml"
	} else {
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GCONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GCONFIG); err != nil {
		fmt.Println(err)
	}
	global.GCONFIG.App.Root, _ = filepath.Abs("..")
	return v
}
