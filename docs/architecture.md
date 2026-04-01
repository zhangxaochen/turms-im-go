# Turms 架构全景分析与 Go 映射策略

## 1. 原始架构概述 (Java 版)

Turms 是一个专为 **10万到1000万级别并发** 设计的开源即时通讯 (IM) 引擎。它的核心设计理念是追求极致性能、高可用性和可伸缩性。

### 1.1 核心组件划分
- **turms-gateway**: 接入层网关。负责维护与海量客户端的 TCP/WebSocket 长连接，管理 Session 生命周期（含多端登录冲突解决、心跳维持），并将请求路由给内网服务。**有状态（内存维护 Session）**。
- **turms-service**: 业务逻辑层。负责处理所有的 IM 业务逻辑（消息投递、群组管理、好友关系、用户资料等）。并直接与存储层 (MongoDB/Redis) 交互。**完全无状态**。
- **turms-server-common**: 公共基础设施层。包含集群发现、自研 RPC 通信、Snowflake ID 生成器、存储层 (Mongo/Redis) 客户端封装、插件框架等。

### 1.2 关键架构决策
- **读扩散模型 (Read-Diffusion)**: 群聊消息只在 MongoDB 存一份。发送消息时，服务直接将消息推送到当前在线的群成员 (Gateway 内存)，不在线的成员上线后通过拉取 (Pull) 方式同步未读消息。
- **纯异步响应式 (Reactive)**: 整个技术栈基于 Project Reactor 和 Netty，从网络 I/O 到 Redis、MongoDB 操作统统是非阻塞的。
- **分层存储 (Tiered Storage)**: MongoDB 数据按照 `deliveryDate` (消息投递时间) 作为 Shard Key 进行分片，天然支持冷热数据分离归档。

---

## 2. Go 重写映射策略

在将其重写为 Go 时，需保持系统宏观架构、存储设计和网络协议的**完全兼容**，但在微观实现上要遵循 Go 的惯用法和高并发模型。

### 2.1 运行时与并发模型映射
| Java (Project Reactor) | Go |
| :--- | :--- |
| `Mono<T>` / `Flux<T>` | `T, error` / `<-chan T` (大部分业务使用同步写法，依靠 `goroutine` 异步) |
| Netty EventLoopGroup | `net.Listener` + goroutine-per-connection (或 goroutine pool) |
| `ConcurrentHashMap` (Session 维护) | `sync.Map` 或带 RWMutex 的分片 `map` 降低锁竞争 |

### 2.2 服务间通信与 RPC
- 原选用自研二进制 RPC 而非 gRPC 是为了省去 gRPC/HTTP2 带来的 overhead。
- **Go 策略**: 依然实现原系统的 custom RPC 帧格式解码器 (`RpcFrameDecoder`/`Encoder`)。可以直接用 net.TCP + 自定义协议处理 goroutine 来实现，保持与 Java 版组件兼容（如果要混合部署的话）；如果纯 Go 部署，可以考虑用高性能的 Go 原生 RPC 或专门优化的 gRPC，但**推荐严格复刻原二进制 RPC** 以实现无缝平滑替换。

### 2.3 存储访问模型
- MongoDB: 从 `spring-data-mongodb-reactive` 替换为官方 `go.mongodb.org/mongo-driver/mongo`。业务层面的读写需要剥离 Reactor，采用标准的 `context.Context` 控制超时和传播。
- Redis: 从 `lettuce` 替换为 `go-redis/redis/v8` (或 v9)。

### 2.4 面向对象映射
- 原来的实体类使用 `@Data` 注解，Go 中映射为带 `bson` 和 `json` tag 的 `struct`。
- Spring 原有的 `@Service` 和自动装配 `@Autowired` 在 Go 中改写为显式的结构体工厂函数（例如 `NewMessageService(mongoClient, redisClient)`），拒绝使用隐式的反转控制(IOC)框架，提升代码可读性和启动性能。
