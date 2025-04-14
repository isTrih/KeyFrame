package notification

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/notification"
	"zerobackend/internal/svc"
)

func MarkAllNotificationReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := notification.NewMarkAllNotificationReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkAllNotificationRead()
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
