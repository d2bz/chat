package websocket

type Message struct {
	Method string      `json:"method"`
	FromId string      `json:"fromId"`
	Data   interface{} `json:"data"` // map[string]interface{}
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FromId: formId,
		Data:   data,
	}
}
