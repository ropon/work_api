package dblayer

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"workApi/config"
)

var redisDb *redis.Client

func InitRedis() (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     config.RHostPort(),
		Password: config.Cfg.Password,
		DB:       config.Cfg.DB,
	})
	_, err = redisDb.Ping().Result()
	return
}

func RedisSet(key string, val interface{}, expiration int64) (err error) {
	err = redisDb.Set(key, val, time.Second*time.Duration(expiration)).Err()
	if err != nil {
		return
	}
	return
}

func RedisGet(key string) (val interface{}, err error) {
	val, err = redisDb.Get(key).Result()
	if err == redis.Nil {
		err = fmt.Errorf("键%s对应值不存在", key)
	} else if err != nil {
		return
	}
	return
}
