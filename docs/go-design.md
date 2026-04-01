# Go 版本总体设计与代码结构 (Go Design)

由于 Java 与 Go 风格差距显著，盲目 1:1 翻译代码将导致在 Go 中产生严重的设计异味 (Code Smells)。以下是我们进行 Turms 架构翻写的模块设计和落地规范。

## 1. 消除框架锁定
- **No Spring**: 抛弃 `@Service`，`@Autowired` 和 `@Value`。依靠显示的结构体定义和简单的 `Options` 函数进行依赖注入 (DI)。
- **配置模型**: 使用 `spf13/viper` 读取配置并 unmarshal 成强类型 `Config` 结构体，不向外暴漏全局配置实例。

## 2. 目录结构规范
参考 Standard Go Project Layout，设计分层和模块切分：
```text
/cmd
  /turms-gateway       (入口: gateway 进程 main)
  /turms-service       (入口: service 进程 main)
/internal
  /cluster             (基础设施: 包含节点管理, 发现)
  /rpc                 (基础设施: 极简二进制 TCP rpc 封装)
  /models              (BSON / Protobuf 模型集中地)
  /storage             (MongoDB & Redis 驱动与常用 Query DAO 封装)
  /domain
    /message           (业务模块划分)
    /user
    /group
    /session
/pkg
  /snowflake           (无状态, 高内聚, 单独封装可测试的算法包)
  /protocol            (protobuf 自动生成的客户端使用包)
/config                (全局参数结构声明)
/docs                  (...)
```

## 3. 依赖注入与服务定义
在 `internal/domain/message/service.go` 中：

```go
// 通过接口约束服务，方便 TDD Mock
type Service interface {
    SendMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.Message, error)
}

type messageService struct {
    mongoCli *mongo.Client
    redisCli *redis.Client
    nodeInfo *cluster.NodeInfo
}

// NewMessageService 作为“构造工厂”传入外部资源
func NewMessageService(m *mongo.Client, r *redis.Client, node *cluster.NodeInfo) Service {
    return &messageService{
        mongoCli: m,
        redisCli: r,
        nodeInfo: node,
    }
}
```

## 4. 针对高并发的内存治理
1. **Zero-Allocation**: 在编解码 RPC 协议包或 WebSocket 帧时使用 `sync.Pool` 借用大块 Byte Buffer 以减轻 GC 负担。
2. **读多写少结构**: Gateway 管理在线网关映射，如果连接字典高达几十上百万，直接使用单一 `sync.RWMutex` + `Map` 会形成灾难级锁竞争。必须实现 **分片 Map (Sharded Concurrent Map)** 减轻冲突（参考 `orcaman/concurrent-map` 的实现原理对 sessionPool 散列分出 N 个 slot）。
