package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"zerobackend/internal/utils"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNewCommentLogic // 创建评论
func NewNewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewCommentLogic {
	return &NewCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewCommentLogic) NewComment(req *types.NewCommentRequest) (resp *types.NewCommentResponse, err error) {
	// 检查违禁词
	insp, _ := utils.DoInsp(l.svcCtx.Config, req.RawContent)

	// 查询ip
	var region = ""
	var ipToUse string
	if req.XIP != "" {
		ipToUse = req.XIP
		parsedIP := net.ParseIP(ipToUse)
		if parsedIP.To4() != nil {
			region, err = l.svcCtx.IP4Searcher.Search(req.XIP)
			if err != nil {
				fmt.Printf("failed to SearchIP(%s): %s\n", req.XIP, err)
			}
		} else if parsedIP.To16() != nil {
			region, err = l.svcCtx.IP6Searcher.Search(req.XIP)
			if err != nil {
				fmt.Printf("failed to SearchIP(%s): %s\n", req.XIP, err)
			}
		}
	} else {
		ipToUse = req.KIP
		parsedIP := net.ParseIP(ipToUse)
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
	}

	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()

	if req.ParentId == 0 {
		// 直接回复
		//TODO:添加通知
		_, newCommentErr := l.svcCtx.CommentModel.ReplyArticle(l.ctx, int64(req.ArticleId), uid, int64(insp), req.Content, region)
		if newCommentErr != nil {
			// 出错
			return nil, newCommentErr
		}
	} else {
		//TODO:添加通知

		_, newCommentErr := l.svcCtx.CommentModel.ReplyComment(l.ctx, int64(req.ArticleId), int64(req.ParentId), int64(req.ParentUserId), uid, int64(insp), req.Content, region)
		if newCommentErr != nil {
			// 出错
			return nil, newCommentErr
		}
	}

	// 构建返回值
	resp = &types.NewCommentResponse{
		Status: "success",
	}
	return resp, nil
}
