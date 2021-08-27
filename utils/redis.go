package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ropon/work_api/conf"
	"time"
)

func RedisSet(key string, val interface{}, expiration int64) (err error) {
	err = conf.RedisCi.Set(key, val, time.Second*time.Duration(expiration)).Err()
	return
}

func RedisGet(key string) (val interface{}, err error) {
	val, err = conf.RedisCi.Get(key).Result()
	if err == redis.Nil {
		err = fmt.Errorf("键%s对应值不存在", key)
	}
	return
}
