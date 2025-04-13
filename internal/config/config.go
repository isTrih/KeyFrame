package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 配置文件中添加一个 Auth 字段，用来配置 JWT 相关的配置
	Auth struct {
		AccessSecret string // 密钥
		AccessExpire int64  // 过期时间
	}
	// MYSQL数据库
	DB struct {
		DataSource string // 数据源
	}
	// PSOTGRESQL数据库
	PG struct {
		DataSource string // 数据源
	}
	// 七牛云储存
	Qiniu struct {
		AK string //
		SK string //
	}
	// UniSMS短信
	Unisms struct {
		SK string
	}
	// IP地址归属地检测
	IPCheck struct {
		Path4 string
		Path6 string
		KEY   string
	}
	// AI接口Key
	Insp struct {
		KEY   string
		URL   string
		MODEL string
	}
	Cache    cache.CacheConf // 配置文件中添加一个 Cache 字段，用来配置Redis缓存相关的配置
	BizRedis redis.RedisConf // 配置文件中添加一个 BizRedis 字段，用来配置业务Redis相关的配置

	// NATS消息队列
	NATS struct {
		ADDR string
	}

	// 邀请码密钥
	InviteKey struct {
		KEY string
		IV  string
	}
}
