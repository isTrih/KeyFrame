package wiki

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zeromicro/x/errors"
	"strconv"
	"time"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/mdl/wiki"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWikiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateWikiLogic // 创建 Wiki
func NewCreateWikiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWikiLogic {
	return &CreateWikiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWikiLogic) CreateWiki(req *types.CreateWikiRequest) (resp *types.StatusResponse, err error) {
	// 检查权限
	utjson, _ := l.ctx.Value("UTYPE").(json.Number)
	ut, _ := utjson.Int64()
	if ut%1000/100 < 5 { // 检查第3位数是否>=5
		return nil, errors.New(7001, "权限不足")
	}

	// 获取用户信息
	uidJson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidJson.Int64()

	tempUser, _ := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	// 构建editor信息
	editorInfo := map[string]interface{}{
		"uid":  uid,
		"name": tempUser.Nickname, // 假设用户名字段为Nickname
		"time": time.Now().Unix(),
	}
	editorJson, _ := json.Marshal(editorInfo)

	// 构建name字段(对应API中的Title)
	nameJson, _ := json.Marshal(map[string]string{"zh": req.Title})

	// 创建新的wiki记录
	wikiData := &wiki.Wiki{
		Type:       sql.NullString{String: strconv.Itoa(int(req.Type)), Valid: true},
		RawContent: sql.NullString{String: req.RawContent, Valid: true},
		Content:    sql.NullString{String: req.Content, Valid: true},
		Editor:     sql.NullString{String: string(editorJson), Valid: true},
		Status:     sql.NullInt64{Int64: 1, Valid: true}, // 默认状态为1(正常)
		Name:       sql.NullString{String: string(nameJson), Valid: true},
	}

	_, err = l.svcCtx.WikiModel.Insert(l.ctx, wikiData)
	if err != nil {
		l.Logger.Errorf("创建Wiki失败: %v", err)
		return nil, errors.New(7002, "创建Wiki失败")
	}

	return &types.StatusResponse{
		Status: "success",
	}, nil
}
