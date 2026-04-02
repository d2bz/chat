package svc

import (
	"chat/apps/im/api/internal/config"
	"chat/apps/im/rpc/im"
	"chat/apps/im/rpc/imclient"
	"chat/apps/social/rpc/social"
	"chat/apps/social/rpc/socialclient"
	"chat/apps/user/rpc/user"
	"chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	im.ImClient
	social.SocialClient
	user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ImClient:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		SocialClient: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		UserClient:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
