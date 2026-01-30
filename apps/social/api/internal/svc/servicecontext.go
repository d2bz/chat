package svc

import (
	"chat/apps/social/api/internal/config"
	"chat/apps/social/api/internal/middleware"
	"chat/apps/social/rpc/socialclient"
	"chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	LimitMiddleware       rest.Middleware
	IdempotenceMiddleware rest.Middleware

	Social socialclient.Social
	User   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		LimitMiddleware:       middleware.NewLimitMiddleware().Handle,
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handle,

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
