# Chat - 即时通讯微服务系统

基于 Go 和 go-zero 框架构建的分布式即时通讯（IM）微服务系统，支持单聊、群聊、好友管理、在线状态等核心功能。

## 技术栈

| 分类 | 技术 |
|------|------|
| 语言 | Go 1.24 |
| 框架 | go-zero |
| 网关 | APISIX |
| 注册中心 | etcd |
| 缓存 | Redis |
| 关系数据库 | MySQL 5.7 |
| 文档数据库 | MongoDB 4.0 |
| 消息队列 | Kafka |
| 链路追踪 | Jaeger |
| 日志分析 | ELK (Elasticsearch + Logstash + Kibana) |
| 配置中心 | Sail |
| 协议 | gRPC、WebSocket、REST |

## 项目架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                              Client                                  │
│                    (Web / Mobile / Desktop)                          │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          APISIX API Gateway                          │
│                     (路由、限流、认证、CORS)                            │
└─────────────────────────────────────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          │                         │                         │
          ▼                         ▼                         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   User Service  │     │  Social Service │     │   IM Service   │
│  (API + RPC)    │     │  (API + RPC)    │     │  (API + RPC)    │
└─────────────────┘     └─────────────────┘     └─────────────────┘
          │                         │                         │
          │                         │                         │
          ▼                         ▼                         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│      MySQL      │     │      MySQL      │     │    MongoDB      │
│  (用户数据)      │     │  (好友、群组)    │     │   (聊天记录)     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    WebSocket Server (im-ws)                          │
│                 实时消息推送、长连接管理                                 │
└─────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Task MQ (Kafka Consumer)                          │
│              消息异步处理、群聊已读合并、会话更新                          │
└─────────────────────────────────────────────────────────────────────┘
```

## 服务模块

| 服务 | 端口 | 说明 |
|------|------|------|
| `user-api` | HTTP | 用户注册、登录、个人信息 |
| `user-rpc` | gRPC | 用户数据查询、RPC 调用 |
| `social-api` | HTTP | 好友申请、群组管理 |
| `social-rpc` | gRPC | 好友关系、群组成员处理 |
| `im-api` | HTTP | 聊天记录、会话管理 |
| `im-rpc` | gRPC | IM 业务逻辑处理 |
| `im-ws` | WebSocket | 实时消息推送、长连接 |
| `task-mq` | Kafka Consumer | 消息持久化、异步推送、已读处理 |

## 目录结构

```
chat/
├── apps/                          # 业务服务
│   ├── user/                      # 用户服务
│   │   ├── api/                  # HTTP API 层
│   │   │   ├── internal/
│   │   │   │   ├── handler/      # 请求处理器
│   │   │   │   ├── logic/        # 业务逻辑
│   │   │   │   ├── svc/          # 服务上下文
│   │   │   │   └── config/       # 配置定义
│   │   │   ├── user.go          # 主入口
│   │   │   └── etc/             # 配置文件
│   │   ├── rpc/                 # gRPC 服务
│   │   └── models/              # 数据模型
│   │
│   ├── social/                   # 社交服务（好友、群组）
│   │   ├── api/
│   │   ├── rpc/
│   │   └── socialmodels/        # 社交领域模型
│   │
│   ├── im/                       # 即时通讯服务
│   │   ├── api/                 # HTTP 接口（聊天记录、会话）
│   │   ├── rpc/                 # gRPC 服务
│   │   ├── ws/                  # WebSocket 服务
│   │   │   ├── websocket/       # WebSocket 核心实现
│   │   │   ├── internal/
│   │   │   │   ├── handler/     # 消息处理器
│   │   │   │   └── svc/
│   │   │   └── ws/              # 消息结构体定义
│   │   └── immodels/            # IM 领域模型
│   │
│   └── task/                     # 任务服务
│       └── mq/                   # Kafka 消费者
│           ├── internal/
│           │   └── handler/
│           │       └── msgTransfer/  # 消息传输处理
│           │           ├── msgChatTransfer.go   # 消息持久化+推送
│           │           ├── msgReadTransfer.go  # 已读处理
│           │           └── groupMsgRead.go     # 群聊已读合并
│           └── mq/               # MQ 客户端
│
├── pkg/                          # 公共包
│   ├── bitmap/                   # 位图（群聊已读标记）
│   ├── configserver/             # 配置中心客户端（Sail）
│   ├── constants/               # 常量定义
│   ├── ctxdata/                 # 上下文数据
│   ├── encrypt/                 # 加密工具
│   ├── interceptor/             # gRPC 拦截器
│   │   ├── idempotence.go       # 幂等性控制
│   │   └── rpcserver/           # RPC 服务端拦截器
│   ├── middleware/              # HTTP 中间件（限流）
│   ├── resultx/                 # 统一响应结构
│   ├── retryjob/                # 重试任务
│   ├── wuid/                    # 唯一 ID 生成
│   └── xerr/                    # 自定义错误
│
├── components/                   # 基础设施配置
│   ├── apisix/                  # APISIX 网关配置
│   ├── apisix-dashboard/        # APISIX Dashboard
│   ├── elasticsearch/           # ES 数据目录
│   ├── etcd/                   # etcd 数据目录
│   ├── kafka/                   # Kafka 数据目录
│   ├── kibana/                 # Kibana 配置
│   ├── logstash/               # Logstash 管道配置
│   ├── mongo/                  # MongoDB 数据目录
│   ├── mysql/                  # MySQL 数据目录
│   ├── redis/                  # Redis 配置
│   └── sail/                   # Sail 配置中心配置
│
├── deploy/                      # 部署相关
│   ├── dockerfile/             # 各服务 Dockerfile
│   ├── mk/                     # Makefile 片段
│   ├── script/                 # 部署脚本
│   └── sql/                   # 数据库建表脚本
│
├── docker-compose.yml           # 基础设施编排
├── go.mod                      # 依赖管理
├── go.sum
└── Makefile                    # 项目构建入口
```

## 核心功能

### 用户模块

- **注册/登录**: 手机号 + 密码认证，JWT Token
- **个人信息**: 获取用户详情

### 社交模块

- **好友管理**: 申请添加好友、审批通过/拒绝、好友列表
- **群组管理**: 创建群组、申请入群、群主/管理员审批、群成员列表
- **在线状态**: 好友在线状态查询

### 即时通讯模块

- **单聊**: 私聊消息发送与接收
- **群聊**: 群组消息广播
- **会话管理**: 会话列表、创建会话、已读未读状态
- **聊天记录**: 分页查询聊天历史

### 消息队列处理

- **消息持久化**: 聊天消息异步写入 MongoDB
- **消息推送**: Kafka 消费后通过 WebSocket 实时推送
- **已读处理**:
  - 单聊已读: 直接更新并推送
  - 群聊已读: **合并机制** — 聚合一定数量或超时后批量推送，避免雪崩
- **会话更新**: 新消息自动更新会话

### 技术亮点

- **WebSocket 长连接**: 自研 WebSocket 服务框架，支持心跳、ACK 确认
- **消息幂等**: 基于 Redis + gRPC 拦截器实现接口幂等
- **限流**: 基于令牌桶的 HTTP/RPC 限流
- **分布式追踪**: Jaeger 全链路日志
- **配置热更新**: Sail 配置中心支持运行时配置刷新 + 优雅重启
- **群聊已读合并**: 位图标记 + 定时/计数批量推送策略

## API 列表

### 用户服务 `/v1/user`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/login` | 用户登录 | 否 |
| POST | `/register` | 用户注册 | 否 |
| GET | `/user` | 获取用户信息 | JWT |

### 社交服务 `/v1/social`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/friend/putIn` | 申请添加好友 | JWT |
| PUT | `/friend/putIn` | 处理好友申请 | JWT |
| GET | `/friend/putIns` | 好友申请列表 | JWT |
| GET | `/friends` | 好友列表 | JWT |
| GET | `/friends/online` | 好友在线状态 | JWT |
| POST | `/group` | 创建群组 | JWT |
| POST | `/group/putIn` | 申请加入群组 | JWT |
| PUT | `/group/putIn` | 处理入群申请 | JWT |
| GET | `/group/putIns` | 入群申请列表 | JWT |
| GET | `/group/users` | 群成员列表 | JWT |
| GET | `/group/users/online` | 群在线用户 | JWT |
| GET | `/groups` | 用户群组列表 | JWT |

### IM 服务 `/v1/im`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/chatlog` | 获取聊天记录 | JWT |
| GET | `/chatlog/readRecords` | 获取已读记录 | JWT |
| GET | `/conversation` | 获取会话列表 | JWT |
| PUT | `/conversation` | 更新会话 | JWT |
| POST | `/setup/conversation` | 创建会话 | JWT |

### WebSocket 消息 (`im-ws`)

| Method | 说明 |
|--------|------|
| `user.online` | 用户上线 |
| `conversation.chat` | 发送聊天消息 |
| `conversation.push` | 消息推送 |
| `conversation.markRead` | 标记消息已读 |

## 快速开始

### 环境要求

- Go 1.24+
- Docker & Docker Compose
- Make

### 1. 启动基础设施

```bash
docker-compose up -d
```

启动 etcd、Redis、MySQL、MongoDB、Zookeeper、Kafka、APISIX、Jaeger、ELK 等服务。

### 2. 初始化数据库

```bash
# MySQL 建表
mysql -h127.0.0.1 -uroot -pchat -e "source deploy/sql/user.sql"
mysql -h127.0.0.1 -uroot -pchat -e "source deploy/sql/social.sql"
mysql -h127.0.0.1 -uroot -pchat -e "source deploy/sql/sail.sql"
```

### 3. 配置 Sail 配置中心

访问 `http://localhost:9000`，创建项目并配置各服务的 YAML 配置。

### 4. 编译并启动服务

```bash
# 编译所有服务
make release-test

# 或单独启动某个服务（开发模式）
make user-api-dev
make im-ws-dev
make task-mq-dev
```

### 5. 通过 APISIX 访问

所有 API 通过 APISIX 网关统一入口:

- 用户服务: `http://localhost:9080/v1/user/`
- 社交服务: `http://localhost:9080/v1/social/`
- IM 服务: `http://localhost:9080/v1/im/`

## 配置说明

各服务的配置文件位于 `apps/{服务名}/etc/dev/` 目录下，支持通过 Sail 配置中心热更新。

主要配置项:

- `Redis`: 会话缓存、在线状态、限流、幂等
- `MySQL`: 用户数据、社交关系
- `MongoDB`: 聊天记录持久化
- `Kafka`: 消息异步处理队列
- `etcd`: 服务注册与发现
- `JWT`: Token 认证密钥

## 数据库表

### MySQL

- `users` — 用户信息
- `friends` — 好友关系
- `friend_requests` — 好友申请
- `groups` — 群组信息
- `group_members` — 群成员
- `group_requests` — 加群申请

### MongoDB

- `chat_log` — 聊天记录集合

## 基础设施端口

| 服务 | 端口 |
|------|------|
| etcd | 3379 |
| Redis | 16379 |
| MySQL | 13306 |
| MongoDB | 47017 |
| Zookeeper | 2181 |
| Kafka | 9092 |
| APISIX Dashboard | 9000 |
| APISIX | 9080 / 9443 |
| Jaeger UI | 16686 |
| Elasticsearch | 9200 |
| Kibana | 5601 |
| Sail Config | 8108 |

## 开发指南

### 新增服务

参考 `apps/user` 的目录结构，使用 goctl 生成:

```bash
goctl api new user-api -path apps/user/api
goctl rpc new user-rpc -path apps/user/rpc -proto xxx.proto
```

### 添加 API 路由

在 `apps/{service}/api/internal/handler/routes.go` 中注册路由。

### 添加 WebSocket 消息类型

1. 在 `apps/im/ws/ws/ws.go` 中定义数据结构
2. 在 `apps/im/ws/internal/handler/router.go` 中注册处理器
3. 在 `apps/task/mq/internal/handler/msgTransfer/` 中添加对应的 MQ 消费处理

### 添加幂等接口

在 `pkg/interceptor/idempotence.go` 的 `method` map 中注册 RPC 方法路径。

## License

MIT
