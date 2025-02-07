package upload

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/upload"
	"zerobackend/internal/svc"
)

// 获取到上传token
func GetUpTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := upload.NewGetUpTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetUpToken()
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
