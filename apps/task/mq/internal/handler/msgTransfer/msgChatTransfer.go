package msgTransfer

import (
	"chat/apps/im/immodels"
	"chat/apps/im/ws/websocket"
	"chat/apps/social/rpc/socialclient"
	"chat/apps/task/mq/internal/svc"
	"chat/apps/task/mq/mq"
	"chat/pkg/constants"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, " value : ", value)
	var (
		data mq.MsgChatTransfer
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据，存入mongodb
	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	// 推送消息
	// 这里是kafka作为ws客户端将消息发给服务端
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = m.single(ctx, &data)
	case constants.GroupChatType:
		err = m.group(ctx, &data)
	}
	return err
}

func (m *MsgChatTransfer) single(ctx context.Context, data *mq.MsgChatTransfer) error {
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		Id:        "",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 查询群的用户
	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, 0, len(users.List))
	//fmt.Printf("group user: %+v", users.List)
	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, members.UserId)
	}
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		Id:        "",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})

}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	// 更新会话
	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
