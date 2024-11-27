package utils

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/x/errors"
	"strconv"
	"time"
)

// RedisCheck 验证码验证
// 示例：
//
//	 rds, rds2 := utils.RedisCheck(req.Mobile, req.VerifyCode)
//		if rds != nil { return nil, rds }
//		if rds2 != nil { return nil, rds2 }
func RedisCheck(mobile string, verifyCode string) (error, error) {
	conf := redis.RedisConf{
		Host:        "106.54.6.216:6379",
		Type:        "node",
		Pass:        "chaozj123123.",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	code, rdserr := rds.Get(mobile)
	if rdserr != nil {
		return nil, errors.New(-1, "缓存获取失败")
	}

	if code != verifyCode && code != "" {
		return nil, errors.New(1001, "验证码错误")
	}
	if code == "" {
		return nil, errors.New(1002, "验证码过期或未获取")
	}

	_, err := rds.Del(mobile)
	if err != nil {
		return nil, errors.New(-1, "缓存删除失败")
	}
	return nil, nil
}

// RedisStorage 验证码存储
// 示例：
//
//		rdsErr := utils.RedisStorage(req.Mobile, 123456, 600)
//	 if rdsErr != nil { return nil, rdsErr }
func RedisStorage(key string, value int, expire int) (err error) {
	conf := redis.RedisConf{
		Host:        "106.54.6.216:6379",
		Type:        "node",
		Pass:        "chaozj123123.",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	err = rds.Setex(key, strconv.Itoa(value), expire)
	return err
}
