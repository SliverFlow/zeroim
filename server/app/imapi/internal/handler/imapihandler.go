package handler

import (
	"net/http"

	"github.com/SliverFlow/zeroim/server/app/imapi/internal/logic"
	"github.com/SliverFlow/zeroim/server/app/imapi/internal/svc"
	"github.com/SliverFlow/zeroim/server/app/imapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImapiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImapiLogic(r.Context(), svcCtx)
		resp, err := l.Imapi(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
