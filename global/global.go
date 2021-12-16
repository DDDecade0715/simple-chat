package global

import (
	"gin-derived/config"
	"gin-derived/pkg/elasticsearch"
	"gin-derived/pkg/rabbitmq"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	GDB       *gorm.DB
	GREDIS    *redis.Client
	GCONFIG   config.Server
	GVIPER    *viper.Viper
	GLOG      *zap.SugaredLogger
	GRABBITMQ *rabbitmq.RabbitMQ
	GES       *elasticsearch.Elasticsearch
)
