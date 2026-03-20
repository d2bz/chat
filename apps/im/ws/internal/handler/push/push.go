package push

import (
	"chat/apps/im/ws/internal/svc"
	"chat/apps/im/ws/websocket"
	"chat/apps/im/ws/ws"
	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
		}

		// 发送给目标
		rconn := srv.GetConn(data.RecvId)
		if rconn == nil {
			// todo: 目标离线
			return
		}

		srv.Infof("push msg %v", data)

		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
		}), rconn)

	}
}
