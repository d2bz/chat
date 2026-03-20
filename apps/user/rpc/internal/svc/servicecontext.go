package svc

import (
	"chat/apps/user/models"
	"chat/apps/user/rpc/internal/config"
	"chat/pkg/constants"
	"chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	UsersModel models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(),
		999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}

	// 写入到redis
	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
