package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zhengjianyang/goCzdb"
	"zerobackend/internal/config"
	"zerobackend/mdl/article"
	"zerobackend/mdl/article_metrics"
	"zerobackend/mdl/comment"
	"zerobackend/mdl/media"
	"zerobackend/mdl/report"
	"zerobackend/mdl/user"
	"zerobackend/mdl/user_action"
	"zerobackend/mdl/user_follow"
	"zerobackend/mdl/wiki"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	BizRedis            *redis.Redis
	Config              config.Config
	IP4Searcher         *goCzdb.DbSearcher
	IP6Searcher         *goCzdb.DbSearcher
	UserModel           user.UserModel
	ArticleModel        article.ArticleModel                // 文章表
	ArticleMetricsModel article_metrics.ArticleMetricsModel // 文章统计表
	MediaModel          media.MediaModel
	UserActionModel     user_action.UserActionModel // 用户行为记录表
	UserFollowModel     user_follow.UserFollowModel // 用户关注关系表
	CommentModel        comment.CommentModel        //评论表
	ReportModel         report.ReportModel          // 举报表
	WikiModel           wiki.WikiModel              //百科表
}

// NewServiceContext 初始化服务上下文
func NewServiceContext(c config.Config) *ServiceContext {

	// 初始化IP4查询器
	searcher4, err4 := goCzdb.NewDbSearcher(c.IPCheck.Path4, "MEMORY", c.IPCheck.KEY)
	if err4 != nil {
		fmt.Printf("Ipv4查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path4, err4)
	}
	// 初始化IP6查询器
	searcher6, err6 := goCzdb.NewDbSearcher(c.IPCheck.Path6, "MEMORY", c.IPCheck.KEY)
	if err6 != nil {
		fmt.Printf("Ipv6查询器 failed to load content from `%s`: %s\n", c.IPCheck.Path6, err4)
	}

	// KeyframeGo := sqlx.NewMysql(c.DB.DataSource)
	keyframeGo := postgres.New(c.PG.DataSource)
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Type: c.BizRedis.Type})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		IP4Searcher: searcher4,
		IP6Searcher: searcher6,
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
		ArticleMetricsModel: article_metrics.NewArticleMetricsModel(keyframeGo, c.Cache),
		// 评论表
		CommentModel: comment.NewCommentModel(keyframeGo, c.Cache),
		// 举报表
		ReportModel: report.NewReportModel(keyframeGo, c.Cache),
		// 百科表
		WikiModel: wiki.NewWikiModel(keyframeGo, c.Cache),
		// 数据库
		BizRedis: rds,
	}
}
