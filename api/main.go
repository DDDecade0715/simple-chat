package api

import (
	"context"
	"fmt"
	"gin-derived/api/routes"
	"gin-derived/global"
	"gin-derived/initialize"
	ws "gin-derived/pkg/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	global.GVIPER = initialize.InitViper()                  //初始化viper
	global.GLOG = initialize.InitLogger("api", "./api/log") //初始化日志
	//global.GREDIS = initialize.InitRedis()                  //初始化redis
	global.GDB = initialize.InitDB() //初始化数据库
	//global.GRABBITMQ = initialize.InitRabbitMQ()            //初始化RabbitMQ
	//global.GES = initialize.InitElasticsearch()             //初始化Elasticsearch
}

func Api() {
	//获取基本配置
	config := global.GCONFIG
	//设置gin模式
	gin.SetMode(config.App.Mode)

	//开启websocket
	go ws.ClientManager.Start()

	//设置http服务
	router := routes.InitRoute()
	var server *http.Server
	go func() {
		server = &http.Server{
			Addr:           ":" + config.App.Port,
			Handler:        router,
			ReadTimeout:    config.App.ReadTimeout * time.Second,
			WriteTimeout:   config.App.WriteTimeout * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			global.GLOG.Fatalf("server.ListenAndServe error: %v", err)
		}
	}()
	//等待中断信号
	quit := make(chan os.Signal)
	//接收syscall.SIGINT和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced shutdown", err)
	}
	fmt.Println("Server existing")
}
