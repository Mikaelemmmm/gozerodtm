package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"gozerodtm/order-api/internal/logic"
	"gozerodtm/order-api/internal/svc"
	"gozerodtm/order-api/internal/types"
)

func createHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuickCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateLogic(r.Context(), ctx)
		resp, err := l.Create(req,r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
