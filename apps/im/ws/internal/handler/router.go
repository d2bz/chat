package handler

import (
	"chat/apps/im/ws/internal/handler/user"
	"chat/apps/im/ws/internal/svc"
	"chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
	})
}
