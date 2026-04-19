package svc

import (
	"chat/apps/im/rpc/imclient"
	"chat/apps/social/api/internal/config"
	"chat/pkg/interceptor"
	"chat/pkg/middleware"

	"chat/apps/social/rpc/socialclient"
	"chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	LimitMiddleware       rest.Middleware
	IdempotenceMiddleware rest.Middleware

	socialclient.Social
	userclient.User
	imclient.Im

	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 10),
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handle,

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient),
		)),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:   imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),

		Redis: redis.MustNewRedis(c.Redisx),
	}
}
