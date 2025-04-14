package notification

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"zerobackend/internal/logic/notification"
	"zerobackend/internal/svc"
)

func GetUnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := notification.NewGetUnreadNotificationCountLogic(r.Context(), svcCtx)
		resp, err := l.GetUnreadNotificationCount()
		if err != nil {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			// code-data 响应格式
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
