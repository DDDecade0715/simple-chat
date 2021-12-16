package initialize

import (
	"gin-derived/global"
	redisClient "gin-derived/pkg/redis"
	"github.com/go-redis/redis"
)

func InitRedis() (r *redis.Client) {
	RedisOptions := &global.GCONFIG.Redis
	p := redisClient.GetRedis("default", RedisOptions)
	r = p.Client
	pong, err := r.Ping().Result()
	if err != nil {
		global.GLOG.Error("redis connect ping failed, err:", err)
	}
	global.GLOG.Infof("redis connect ping response: %s", pong)
	return
}
