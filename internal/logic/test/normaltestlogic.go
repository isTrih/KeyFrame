package test

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/internal/utils"
)

type NormalTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewNormalTestLogic // 创建帧（文章）
func NewNormalTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NormalTestLogic {
	return &NormalTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NormalTestLogic) NormalTest(req *types.TestRequest) (resp *types.TestResponse, err error) {
	result := utils.DecryptTriDESToNumber("tmh&1101", "hhh&1101", req.Content)
	fmt.Print(result)
	resp = &types.TestResponse{}
	resp.Content = strconv.FormatInt(result, 10)
	return resp, nil
}
