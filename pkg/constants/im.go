package constants

type MType int

const (
	TextMType MType = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota
	SingleChatType
)
