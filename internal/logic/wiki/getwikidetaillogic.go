package wiki

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/x/errors"
	"strconv"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWikiDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetWikiDetailLogic // 获取 Wiki 详情
func NewGetWikiDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWikiDetailLogic {
	return &GetWikiDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWikiDetailLogic) GetWikiDetail(req *types.GetWikiDetailRequest) (resp *types.GetWikiDetailResponse, err error) {
	wikiId := int64(req.Id)
	wikiData, err := l.svcCtx.WikiModel.FindOne(l.ctx, wikiId)
	if err != nil {
		l.Logger.Errorf("获取Wiki详情失败, id: %d, error: %v", wikiId, err)
		return nil, errors.New(7003, "获取Wiki详情失败")
	}

	// 解析name字段为标题
	var nameMap map[string]string
	var title string
	if wikiData.Name.Valid {
		if err := json.Unmarshal([]byte(wikiData.Name.String), &nameMap); err == nil {
			title = nameMap["zh"]
		}
	}

	// 解析类型
	wikiType := uint8(0)
	if wikiData.Type.Valid {
		if typeVal, err := strconv.Atoi(wikiData.Type.String); err == nil {
			wikiType = uint8(typeVal)
		}
	}

	// 解析状态
	status := int16(0)
	if wikiData.Status.Valid {
		status = int16(wikiData.Status.Int64)
	}

	// 转换创建时间和更新时间为时间戳
	createTime := uint64(wikiData.CreateTime.Unix())
	updateTime := uint64(wikiData.UpdateTime.Unix())

	return &types.GetWikiDetailResponse{
		Wiki: types.WikiDetail{
			Id:         uint64(wikiData.Id),
			Title:      title,
			Content:    wikiData.Content.String,
			RawContent: wikiData.RawContent.String,
			Type:       wikiType,
			Editor:     wikiData.Editor.String,
			Status:     status,
			CreateTime: createTime,
			UpdateTime: updateTime,
		},
	}, nil
}
