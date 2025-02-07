package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zhengjianyang/goCzdb"
	"zerobackend/internal/config"
	"zerobackend/mdl/action_count"
	"zerobackend/mdl/article"
	"zerobackend/mdl/media"
	"zerobackend/mdl/user"
	"zerobackend/mdl/user_action"
	"zerobackend/mdl/user_follow"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	BizRedis         *redis.Redis
	Config           config.Config
	IP4Searcher      *goCzdb.DbSearcher
	UserModel        user.UserModel
	ArticleModel     article.ArticleModel
	MediaModel       media.MediaModel
	UserActionModel  user_action.UserActionModel   // 用户行为记录表
	UserFollowModel  user_follow.UserFollowModel   // 用户关注关系表
	ActionCountModel action_count.ActionCountModel // 行为统计表
}

// NewServiceContext 初始化服务上下文
func NewServiceContext(c config.Config) *ServiceContext {

	// 初始化IP查询器
	searcher4, err4 := goCzdb.NewDbSearcher(c.IPCheck.Path4, "MEMORY", c.IPCheck.KEY)
	if err4 != nil {
		fmt.Printf("Ipv4查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path4, err4)
	}

	keyframeGo := sqlx.NewMysql(c.DB.DataSource)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Type: c.BizRedis.Type})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		IP4Searcher: searcher4,
		Config:      c,
		// 用户数据库
		UserModel: user.NewUserModel(keyframeGo, c.Cache),
		// 文章数据库
		ArticleModel: article.NewArticleModel(keyframeGo, c.Cache),
		// 点赞 收藏 关注数据库
		UserActionModel: user_action.NewUserActionModel(keyframeGo, c.Cache),
		// 用户关注关系表
		UserFollowModel: user_follow.NewUserFollowModel(keyframeGo, c.Cache),
		// 用户行为记录表
		ActionCountModel: action_count.NewActionCountModel(keyframeGo, c.Cache),
		BizRedis:         rds,
	}
}
