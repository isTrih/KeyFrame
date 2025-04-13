package wiki

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/wiki"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
)

// GetLatestWikisHandler // 获取最新 Wiki 列表
func GetLatestWikisHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLatestWikisRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := wiki.NewGetLatestWikisLogic(r.Context(), svcCtx)
		resp, err := l.GetLatestWikis(&req)
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
