package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/x/errors"
	"strconv"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
	"zerobackend/mdl/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserInfoLogic 获取用户信息 需要鉴权
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	uidjson, _ := l.ctx.Value("UID").(json.Number)
	typejson, _ := l.ctx.Value("UTYPE").(json.Number)
	uid, _ := uidjson.Int64()
	utype, _ := typejson.Int64()

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")
	}

	resp = new(types.UserInfoResponse)
	resp.UserId = int64(user.Id)
	resp.Username = user.Nickname
	resp.Avatar = user.Mobile
	resp.Type = strconv.FormatInt(utype, 10)
	return
}
