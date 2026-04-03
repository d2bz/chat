package svc

import (
	"chat/apps/im/api/internal/config"
	"chat/apps/im/rpc/imclient"
	"chat/apps/social/rpc/socialclient"
	"chat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	imclient.Im
	socialclient.Social
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
