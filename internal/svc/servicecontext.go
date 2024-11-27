package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zerobackend/internal/config"
	"zerobackend/mdl/user/model"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

// NewServiceContext 初始化服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
