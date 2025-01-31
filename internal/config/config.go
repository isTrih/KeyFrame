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
	DB struct {
		DataSource string // 数据源
	}
	Qiniu struct {
		AK string //
		SK string //
	}
	Unisms struct {
		SK string
	}
	IPCheck struct {
		Path4 string
		KEY   string
	}
	Cache    cache.CacheConf // 配置文件中添加一个 Cache 字段，用来配置Redis缓存相关的配置
	BizRedis redis.RedisConf
}
