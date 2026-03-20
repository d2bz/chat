package ws

import "chat/pkg/constants"

type (
	// 在websocket中，原始数据会被反序列化为map[string]interface{}
	// 因此要引入mapstructure库把map的数据绑定到结构体上
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
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
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`

		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
)
