package comment

import (
	"context"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCommentListLogic // 获取评论列表
func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentListLogic) GetCommentList(req *types.CommentListRequest) (resp *types.CommentListResponse, err error) {
	// 获取原始评论数据
	comments, total, err := l.svcCtx.CommentModel.GetCommentList(l.ctx, req.ArticleId, req.Page)
	if err != nil {
		return nil, err
	}

	// 组织评论结构
	mainComments := make(map[uint64]*types.CommentItem)
	var commentList []*types.CommentItem // 修改为指针类型的切片

	// 先处理主评论
	for _, comment := range comments {
		if comment.ParentId == 0 {
			item := &types.CommentItem{ // 使用指针
				Id:              comment.Id,
				UserId:          comment.UserId,
				Nickname:        comment.Nickname,
				Avatar:          comment.Avatar,
				Content:         comment.Content,
				LikeCount:       comment.LikeCount,
				CreateTime:      comment.CreateTime,
				IpLocation:      comment.IpLocation,
				SubComments:     make([]types.SubCommentItem, 0),
				SubCommentCount: 0,
			}
			mainComments[comment.Id] = item
			commentList = append(commentList, item)
		}
	}

	// 再处理子评论
	for _, comment := range comments {
		if comment.ParentId != 0 {
			if parent, ok := mainComments[comment.ParentId]; ok {
				subComment := types.SubCommentItem{
					Id:              comment.Id,
					UserId:          comment.UserId,
					Nickname:        comment.Nickname,
					Avatar:          comment.Avatar,
					Content:         comment.Content,
					LikeCount:       comment.LikeCount,
					CreateTime:      comment.CreateTime,
					ParentId:        comment.ParentId,
					ParentUserId:    comment.ParentUserId,
					ReplyToNickname: comment.ReplyToNickname,
					IpLocation:      comment.IpLocation,
				}
				parent.SubComments = append(parent.SubComments, subComment)
				parent.SubCommentCount = uint64(len(parent.SubComments))
			}
		}
	}

	// 转换回非指针类型的结果
	result := make([]types.CommentItem, 0, len(commentList))
	for _, item := range commentList {
		result = append(result, *item)
	}

	return &types.CommentListResponse{
		Comments: result,
		Total:    total,
	}, nil
}
