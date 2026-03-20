package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error

	Send(v any) error
	Read(v any) error
}

type client struct {
	*websocket.Conn // 在这个连接对象中本身就有关闭的方法
	host            string

	opt dailOption
}

func NewClient(host string, opts ...DailOptions) Client {
	opt := newDailOptions(opts...)

	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return &c
}

// 建立与websocket服务的连接
// 此处参照websocket/example/echo/client.go中的实现方法
func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}
	// 有错误时才重新创建连接
	// 如果建立连接后，服务端重启，那么客户端拿到的连接就会失效，拿失效的连接发消息会报错
	// 所以尝试重新建立连接
	conn, err := c.dail()
	if err != nil {
		return err
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
