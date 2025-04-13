package wiki

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/wiki"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
)

// CreateWikiHandler // 创建 Wiki
func CreateWikiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateWikiRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := wiki.NewCreateWikiLogic(r.Context(), svcCtx)
		resp, err := l.CreateWiki(&req)
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
