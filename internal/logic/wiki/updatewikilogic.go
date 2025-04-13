package wiki

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type UpdateWikiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateWikiLogic // 编辑 Wiki
func NewUpdateWikiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWikiLogic {
	return &UpdateWikiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWikiLogic) UpdateWiki(req *types.UpdateWikiRequest) (resp *types.StatusResponse, err error) {
	// 检查权限
	utjson, _ := l.ctx.Value("UTYPE").(json.Number)
	ut, _ := utjson.Int64()
	if ut%1000/100 < 5 { // 检查第3位数是否>=5
		return nil, errors.New(6001, "权限不足")
	}

	// 获取用户信息
	uidJson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidJson.Int64()

	// 使用UserModel获取用户信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, errors.New(6005, "获取用户信息失败")
	}
	unameStr := userInfo.Nickname // 假设用户名字段为Nickname

	// 查询原Wiki记录
	wikiId := int64(req.Id)
	oldWiki, err := l.svcCtx.WikiModel.FindOne(l.ctx, wikiId)
	if err != nil {
		l.Logger.Errorf("Wiki不存在, id: %d, error: %v", wikiId, err)
		return nil, errors.New(6002, "wiki不存在")
	}

	// 保存到历史记录表
	err = l.svcCtx.WikiModel.InsertWikiHistory(l.ctx, wikiId, uid, oldWiki.RawContent.String, req.ChangeLog)
	if err != nil {
		l.Logger.Errorf("保存Wiki历史记录失败: %v", err)
		return nil, errors.New(6003, "保存历史记录失败")
	}

	// 构建editor信息
	editorInfo := map[string]interface{}{
		"uid":  uid,
		"name": unameStr,
		"time": time.Now().Unix(),
	}
	editorJson, _ := json.Marshal(editorInfo)

	// 构建name字段(对应API中的Title)
	nameJson, _ := json.Marshal(map[string]string{"zh": req.Title})

	// 更新Wiki
	oldWiki.Type = sql.NullString{String: strconv.Itoa(int(req.Type)), Valid: true}
	oldWiki.RawContent = sql.NullString{String: req.RawContent, Valid: true}
	oldWiki.Content = sql.NullString{String: req.Content, Valid: true}
	oldWiki.Editor = sql.NullString{String: string(editorJson), Valid: true}
	oldWiki.Name = sql.NullString{String: string(nameJson), Valid: true}

	err = l.svcCtx.WikiModel.Update(l.ctx, oldWiki)
	if err != nil {
		l.Logger.Errorf("更新Wiki失败: %v", err)
		return nil, errors.New(6004, "更新Wiki失败")
	}

	return &types.StatusResponse{
		Status: "success",
	}, nil
}
