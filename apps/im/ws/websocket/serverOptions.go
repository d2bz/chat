package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	pattern string

	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOption {
	// 默认配置
	o := serverOption{
		Authentication:    new(authentication),
		pattern:           "/ws",
		maxConnectionIdle: defaultMaxConnectionIdle,
	}
	// 加载自定义配置
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithAuthentication(authentication Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = authentication
	}
}

func WithHandlerPattern(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnection time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnection > 0 {
			opt.maxConnectionIdle = maxConnection
		}
	}
}
