package conversation

import (
	"chat/apps/im/ws/internal/svc"
	"chat/apps/im/ws/websocket"
	"chat/apps/im/ws/ws"
	"chat/apps/task/mq/mq"
	"chat/pkg/constants"
	"github.com/mitchellh/mapstructure"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Chat
		// 对map[string]interface{}类型进行转换
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				SendTime:       time.Now().UnixMilli(),
				MType:          data.MType,
				Content:        data.Content,
			})
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}
		}
	}
}
