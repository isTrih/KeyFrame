package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zhengjianyang/goCzdb"
	"zerobackend/internal/config"
	"zerobackend/mdl/article"
	"zerobackend/mdl/collect"
	"zerobackend/mdl/follow"
	"zerobackend/mdl/follow_count"
	"zerobackend/mdl/like_count"
	"zerobackend/mdl/like_record"
	"zerobackend/mdl/media"
	"zerobackend/mdl/reply"
	"zerobackend/mdl/reply_count"
	"zerobackend/mdl/tag"
	"zerobackend/mdl/tag_resource"
	"zerobackend/mdl/user"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	BizRedis         *redis.Redis
	Config           config.Config
	IP4Searcher      *goCzdb.DbSearcher
	UserModel        user.UserModel
	ArticleModel     article.ArticleModel
	ReplyModel       reply.ReplyModel
	ReplyCountModel  reply_count.ReplyCountModel
	TagModel         tag.TagModel
	TagResourceModel tag_resource.TagResourceModel
	LikeCountModel   like_count.LikeCountModel
	LikeRecordModel  like_record.LikeRecordModel
	FollowModel      follow.FollowModel
	FollowCountModel follow_count.FollowCountModel
	MediaModel       media.MediaModel
	CollectModel     collect.CollectModel
}

// NewServiceContext 初始化服务上下文
func NewServiceContext(c config.Config) *ServiceContext {

	searcher4, err4 := goCzdb.NewDbSearcher(c.IPCheck.Path4, "MEMORY", c.IPCheck.KEY)
	if err4 != nil {
		fmt.Printf("Ipv4查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path4, err4)
	}

	keyframeGo := sqlx.NewMysql(c.DB.DataSource)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		IP4Searcher: searcher4,
		Config:      c,
		//用户数据库
		UserModel: user.NewUserModel(keyframeGo, c.Cache),
		//文章数据库
		ArticleModel: article.NewArticleModel(keyframeGo, c.Cache),
		//关注数量
		FollowCountModel: follow_count.NewFollowCountModel(keyframeGo, c.Cache),
		//收藏
		CollectModel: collect.NewCollectModel(keyframeGo, c.Cache),
		//点赞
		LikeRecordModel: like_record.NewLikeRecordModel(keyframeGo, c.Cache),
		LikeCountModel:  like_count.NewLikeCountModel(keyframeGo, c.Cache),
		BizRedis:        rds,
	}
}
