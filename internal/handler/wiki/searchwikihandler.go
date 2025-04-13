package wiki

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/wiki"
	"zerobackend/internal/svc"
	"zerobackend/internal/types"
)

// SearchWikiHandler // 搜索 Wiki
func SearchWikiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchWikiRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := wiki.NewSearchWikiLogic(r.Context(), svcCtx)
		resp, err := l.SearchWiki(&req)
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
