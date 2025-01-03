package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/x/errors"
	model "zerobackend/mdl/user"

	"zerobackend/internal/svc"
	"zerobackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserInfoLogic // 获取用户信息 不需要token
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	//创建返回数据
	resp = new(types.UserInfoResponse)

	//查询用户基本信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if user == nil {
		return nil, errors.New(6021, "用户不存在")
	}

	resp.UserId = user.Id
	resp.Username = user.Nickname
	resp.Avatar = user.Avatar
	resp.Type = uint8(user.Type)
	resp.Status = uint8(user.Status)
	resp.VNote = user.Vnote
	resp.Signature = user.Signature

	//查询用户的文章数
	feedCount, err := l.svcCtx.ArticleModel.GetFeedsNum(l.ctx, req.UserId)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}

	resp.FeedCount = uint64(feedCount)

	//查询用户的粉丝数
	followCount, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, req.UserId)
	if err != nil && err != model.ErrNotFound {
		fmt.Println(err)
		return nil, errors.New(4001, "查询数据失败")
	}
	if followCount == nil {
		resp.FansCount = 0
		resp.FollowCount = 0
	} else {
		resp.FansCount = followCount.FansCount
		resp.FollowCount = followCount.FollowCount
	}

	//返回数据
	return resp, nil
}
