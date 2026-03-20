package constants

// 消息类型
type MType int

const (
	// 文本类型
	TextMType MType = iota
)

// 聊天类型
type ChatType int

const (
	// 群聊
	GroupChatType ChatType = iota
	// 私聊
	SingleChatType
)
