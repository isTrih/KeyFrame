package utils

import (
	"strconv"
	"time"
	"zerobackend/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/x/errors"
)

// RedisCheck 验证码验证
// 示例：
//
//	 rds, rds2 := utils.RedisCheck(req.Mobile, req.VerifyCode)
//		if rds != nil { return nil, rds }
//		if rds2 != nil { return nil, rds2 }
func RedisCheck(c config.Config, mobile string, verifyCode string) (error, error) {
	conf := redis.RedisConf{
		Host:        c.BizRedis.Host,
		Type:        "node",
		Pass:        c.BizRedis.Pass,
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
func RedisStorage(c config.Config, key string, value int, expire int) (err error) {
	conf := redis.RedisConf{
		Host:        c.BizRedis.Host,
		Type:        "node",
		Pass:        c.BizRedis.Pass,
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	err = rds.Setex(key, strconv.Itoa(value), expire)
	return err
}

func RedisKey(c config.Config, key string) (keys []string, err error) {
	rds := redis.MustNewRedis(c.BizRedis)

	keys, err = rds.Keys(key)

	return keys, err
}
