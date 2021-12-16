package initialize

import (
	"gin-derived/global"
	"gin-derived/pkg/rabbitmq"
)

func InitRabbitMQ() *rabbitmq.RabbitMQ {
	//实例一个rabbitmq
	config := global.GCONFIG.Amqp.Config.Channels.Default
	url := global.GCONFIG.Amqp.Url
	connection := rabbitmq.GetRabbitMQ("default", config, url)

	return connection
}
