package report

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	model "zerobackend/mdl/report"
)

type ReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewReportLogic // 举报
func NewReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportLogic {
	return &ReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportLogic) Report(req *types.ReportRequest) (resp *types.ReportResponse, err error) {

	// 获取用户ID
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	uid, _ := uidjson.Int64()
	_, err = l.svcCtx.ReportModel.Insert(l.ctx, &model.Report{
		ReportedId:     int64(req.ObjectId),
		ReportedType:   int64(req.ObjectType),
		ReporterUserId: uid,
		ReportType:     int64(req.Type),
		OwnerId:        int64(req.OwnerId),
	})
	{

	}
	return
}
