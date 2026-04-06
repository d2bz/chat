package ws

import "chat/pkg/constants"

type (
	// 在websocket中，原始数据会被反序列化为map[string]interface{}
	// 因此要引入mapstructure库把map的数据绑定到结构体上
	Msg struct {
		MsgId           string `mapstructure:"msgId"`
		constants.MType `mapstructure:"mType"`
		Content         string            `mapstructure:"content"`
		ReadRecords     map[string]string `mapstructure:"readRecords"`
	}

	// 对应message中的data
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	// 用于定义接收客户端传入的消息体结构
	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string   `mapstructure:"sendId"`
		RecvId             string   `mapstructure:"recvId"`
		RecvIds            []string `mapstructure:"recvIds"` // 群聊消息推送时考虑多个用户的接收
		SendTime           int64    `mapstructure:"sendTime"`

		MsgId       string                `mapstructure:"msgId"`
		ReadRecords map[string]string     `mapstructure:"readRecords"`
		ContentType constants.ContentType `mapstructure:"contentType"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`
		ConversationId     string   `mapstructure:"conversationId"`
		RecvId             string   `mapstructure:"recvId"`
		MsgIds             []string `mapstructure:"msgIds"`
	}
)
