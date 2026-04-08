package constants

const (
	// REDIS_SYSTEM_ROOT_TOKEN 全局权限用户, 用于mq服务作为客户端向ws服务发送消息的身份凭证
	REDIS_SYSTEM_ROOT_TOKEN string = "system:root:token"
	// REDIS_ONLINE_USER 在线用户，
	REDIS_ONLINE_USER = "online:user"
)
