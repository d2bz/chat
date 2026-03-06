package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	service.ServiceConf // 引入go-zero中的日志处理等功能

	ListenOn string
}
