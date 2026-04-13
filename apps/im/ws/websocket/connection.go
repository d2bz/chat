package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex
	Uid    string
	*websocket.Conn
	s *Server

	messageMu sync.Mutex
	// 读消息队列
	readMessage []*Message
	// 记录消息序列化(key：消息id value：具体消息)
	readMessageSeq map[string]*Message
	// 该通道用于ACK确认之后将消息发送给任务处理(handlerWrite)
	message chan *Message

	idle              time.Time //最近一次进入“空闲状态”的时间
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	var responseHeader http.Header
	if protocol := r.Header.Get("Sec-WebSocket-Protocol"); protocol != "" {
		responseHeader = http.Header{"Sec-WebSocket-Protocol": []string{protocol}}
	}

	c, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:           c,
		s:              s,
		readMessage:    make([]*Message, 0, 2),
		readMessageSeq: make(map[string]*Message, 2),
		// 通道大小设置为1：可以减少数据投递中的阻塞情况，也可以保障这个收和发的执行的顺序问题
		message:           make(chan *Message, 1),
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}),
	}

	go conn.keepalive()
	return conn
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	// 读队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 已经有消息的记录，该消息已经有ack的确认
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}

		// ACK的确认是对当前序号进行+1处理
		// 如果传入的msg序号 <= 队列中存在的消息序号，则说明未进行ACK确认
		if msg.AckSeq <= m.AckSeq {
			// 没有进行ack的确认, 也可能为重复发送
			return
		}

		c.readMessageSeq[msg.Id] = msg
		return
	}
	// 还没有进行ack的确认, 避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

// ReadMessage 和 WriteMessage 是对第三方库的方法重写，加入对空闲时间的状态处理
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 此处第三方库中的方法并不安全，所以需要加锁
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

// 重写Close避免通道被多次关闭
func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // 连接不在空闲状态，重置计时器
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			// 计算还剩多少空闲时间，如果超时就优雅断开连接，未超时就按剩余设置计时器
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		// 连接关闭，退出keepalive
		case <-c.done:
			return
		}
	}
}
