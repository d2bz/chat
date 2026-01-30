package group

import (
	"net/http"

	"chat/apps/social/api/internal/logic/group"
	"chat/apps/social/api/internal/svc"
	"chat/apps/social/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 成员列表列表
func GroupUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupUserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupUserListLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
