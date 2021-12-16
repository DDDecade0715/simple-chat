package redis

import (
	"gin-derived/config"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

var RedisPool = make(map[string]*Redis)

func GetRedis(name string, redisOptions *config.Redis) (r *Redis) {
	var ok bool
	if r, ok = RedisPool[name]; ok {
		return r
	}
	r = CreateRedis(redisOptions.Addr, redisOptions.Password, redisOptions.DB)

	RedisPool[name] = r

	return
}

func CreateRedis(address string, password string, db int) (r *Redis) {
	options := &redis.Options{
		Addr:     address,  // 要连接的redis IP:port
		Password: password, // redis 密码
		DB:       db,       // 要连接的redis 库
	}
	client := redis.NewClient(options)

	r = &Redis{
		Client: client,
	}
	return
}
