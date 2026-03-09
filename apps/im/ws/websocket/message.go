package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0 // 用户消息
	FramePing FrameType = 0x1 // 心跳消息
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method"`
	FromId    string      `json:"fromId"`
	Data      interface{} `json:"data"` // map[string]interface{}
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    formId,
		Data:      data,
	}
}
