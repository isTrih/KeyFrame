package user

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/user"
	"zerobackend/internal/svc"
)

// 获取用户关注\收藏\点赞列表
func UserRelationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := user.NewUserRelationLogic(r.Context(), svcCtx)
		resp, err := l.UserRelation()
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
