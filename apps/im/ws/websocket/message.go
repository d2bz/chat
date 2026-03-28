package websocket

import "time"

type FrameType uint8

const (
	FrameData  FrameType = 0x0 // 用户消息
	FramePing  FrameType = 0x1 // 心跳消息
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9 // 错误类型
)

type Message struct {
	FrameType `json:"frameType"`
	Id        string      `json:"id"`
	AckSeq    int         `json:"ackSeq"`
	ackTime   time.Time   `json:"ackTime"`
	errCount  int         `json:"errCount"`
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

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
