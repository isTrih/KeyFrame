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
	IP6Searcher      *goCzdb.DbSearcher
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
	searcher6, err6 := goCzdb.NewDbSearcher(c.IPCheck.Path6, "MEMORY", c.IPCheck.KEY)

	if err4 != nil {
		fmt.Printf("Ipv4查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path4, err4)
	}
	if err6 != nil {
		fmt.Printf("IPv6查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path6, err6)
	}

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		IP4Searcher: searcher4,
		IP6Searcher: searcher6,
		Config:      c,
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
