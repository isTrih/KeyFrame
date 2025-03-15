package feed

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/feed"
	"zerobackend/internal/svc"
)

// ShareFeedXHSHandler // 分享小红书帧（文章）
func ShareFeedXHSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := feed.NewShareFeedXHSLogic(r.Context(), svcCtx)
		resp, err := l.ShareFeedXHS()
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
