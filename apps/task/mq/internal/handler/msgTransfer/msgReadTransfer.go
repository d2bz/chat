package msgTransfer

import (
	"chat/apps/im/ws/ws"
	"chat/apps/task/mq/internal/svc"
	"chat/apps/task/mq/mq"
	"chat/pkg/bitmap"
	"chat/pkg/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

type MsgReadTransfer struct {
	*baseMsgTransfer
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Infof("MsgReadTransfer ", value)
	var data mq.MsgMarkRead
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理--更新消息为已读
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	})
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	result := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读消息，遍历每一条消息设置当前用户为已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 设置当前发送者用户为已读状态
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		// 将已读消息（二进制）进行base64编码转换，这样可以保证在网络传输过程中的一个精度
		// 前端也可以使用base64进行解码，将已读消息还原为二进制
		result[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err := m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
