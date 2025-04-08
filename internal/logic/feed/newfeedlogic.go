package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/x/errors"
	"net"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNewFeedLogic // 创建帧（文章）
func NewNewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewFeedLogic {
	return &NewFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewFeedLogic) NewFeed(req *types.NewFeedRequest) (resp *types.StatusResponse, err error) {
	// 流程：客户端上传图片（带有自动删除的生命周期），完成发布后更改生命周期
	fmt.Println(req)
	// 更新图片资源有效期为永久
	accessKey := l.svcCtx.Config.Qiniu.AK
	secretKey := l.svcCtx.Config.Qiniu.SK
	bucket := "chaozj-keyframe"
	if req.Cover != "xx" && req.Media != nil {
		threading.GoSafe(func() {
			err := utils.BatchUpdateFileLifecycle(req.Media, req.Cover, accessKey, secretKey, bucket)
			if err != nil {
				fmt.Println("七牛出错", err)
			}
		})
	}
	if req.Cover == "xx" && req.Media != nil {
		threading.GoSafe(func() {
			err := utils.BatchUpdateFileLifecycle(req.Media, "UNCOVER", accessKey, secretKey, bucket)
			if err != nil {
				fmt.Println("七牛出错", err)
			}
		})
	}
	// 检查违禁词
	insp, _ := utils.DoInsp(l.svcCtx.Config, req.RawContent+req.Title)

	// 查询ip
	var region = ""
	parsedIP := net.ParseIP(req.KIP)
	if parsedIP.To4() != nil {
		region, err = l.svcCtx.IP4Searcher.Search(req.KIP)
		if err != nil {
			fmt.Printf("failed to SearchIP(%s): %s\n", req.KIP, err)
		}
	} else if parsedIP.To16() != nil {
		region, err = l.svcCtx.IP6Searcher.Search(req.KIP)
		if err != nil {
			fmt.Printf("failed to SearchIP(%s): %s\n", req.KIP, err)
		}
	}

	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()

	// 发布文章

	if req.Media != nil && req.Cover != "xx" {
		err = l.svcCtx.ArticleModel.NewFeed(
			l.ctx,
			req.Title, req.Content, req.RawContent,
			req.Cover, req.CoverInfo.Height, req.CoverInfo.Width, req.Media,
			uid, region, insp)
	} else {
		err = l.svcCtx.ArticleModel.NewFeed(
			l.ctx,
			req.Title, req.Content, req.RawContent,
			req.Cover, 100, 200, req.Media,
			uid, region, insp)
	}

	//_, err = l.svcCtx.ArticleModel.Insert(l.ctx,
	//	&article.Article{
	//		Title:       req.Title,
	//		Content:     req.Content,
	//		RawContent:  req.RawContent,
	//		AuthorId:    uint64(uid),
	//		Status:      0,
	//		IpLocation:  region,
	//		PublishTime: time.Time{},
	//		CreateTime:  time.Time{},
	//		UpdateTime:  time.Time{},
	//		AiInsp:      uint64(insp),
	//	})
	if err != nil {
		return nil, errors.New(4003, "文章发布失败,请稍后再试")
	}
	resp = new(types.StatusResponse)
	resp.Status = "SUCCESS"
	return resp, nil
}
