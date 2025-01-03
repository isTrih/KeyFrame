package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		//用户数据库
		UserModel: user.NewUserModel(sqlConn, c.Cache),
		//文章数据库
		ArticleModel: article.NewArticleModel(sqlConn, c.Cache),
		//关注数量
		FollowCountModel: follow_count.NewFollowCountModel(sqlConn, c.Cache),
		//收藏
		CollectModel: collect.NewCollectModel(sqlConn, c.Cache),
		//点赞
		LikeRecordModel: like_record.NewLikeRecordModel(sqlConn, c.Cache),
		LikeCountModel:  like_count.NewLikeCountModel(sqlConn, c.Cache),
		BizRedis:        rds,
	}
}
