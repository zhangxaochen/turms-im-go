# 基础设施组件映射设计 (Infrastructure)

Turms 的高性能部分归功于它强大的 `turms-server-common` 基础设施包。在 Go 重写中，这是工作量最大且最核心的基石。

## 1. 集群与节点管理 (`im.turms.server.common.infra.cluster.node.Node`)
- **功能**: 处理服务注册、心跳、领导者选举、以及感知网关/服务节点的增删。
- **Go 替代版设计**:
  - 原理不依赖特定的 CP 协调系统（例如 Zookeeper），可以通过 MongoDB 或 Redis 执行 lease 租约。
  - 需要在 Go 中维护一个并发安全的节点列表 `NodeList` 管理当前在线节点的 IP:RPC-Port 映射表。

## 2. 自研 RPC 框架
Turms 没有选用 gRPC，而是使用长连接和极其精简的二进制帧编解码。
- `RpcRequestCodec` / `RpcResponseCodec` / `RpcFrameDecoder`
- **Go 实现思路**:
  - 基于原生的 `net.Listen("tcp")`。
  - 读取 4 字节的 header 校验包长，避免粘包拆包，直接读取 protobuf 字节流序列化。
  - 连接池复用：跨节点不建议频繁建立连接，建议使用长连接池（如 `fatih/pool` 或自建多路复用 Multiplexer）。

## 3. ID 生成器 (Snowflake)
全局唯一并且大致时间递增的分布式发号器。
- **Go 替代版设计**:
  - 标准的 Twitter 雪花算法 Go 实现 (例如 `bwmarrin/snowflake` 或者直接复刻基于时钟回拨保护的自研版本)。
  - `NodeId` 的分配依赖 MongoDB/Redis 集群分配当前进程持有的唯一 MachineID。

## 4. 插件与扩展机制 (Plugins)
原系统通过 `JavaPlugin` 和基于 GraalVM 的 `JsPlugin` 支持极强的自定义。对于消息过滤、用户权限扩展非常重要。
- **Go 替代版设计**:
  - **Go 原生插件** (通过 `plugin` 标准库，但局限性较大，仅限 Linux，不太推荐)。
  - **Lua / JS 脚本化**：更推荐引入轻量级脚本，如 `dop251/goja` 执行 JavaScript 逻辑，或者 `yuin/gopher-lua` 允许动态注入简单的业务逻辑拦截。
  - **WebAssembly (WASM)**：使用 `wazero` 提供安全无依赖的多语言插件系统（终极且最推荐的解法，支持 Rust/C/Go 等编译的插件拦截）。

## 5. 指标与遥测 (Metrics)
监控大厅性能和内存状态指标。
- **Java**: Micrometer + Prometheus。
- **Go**: 官方的 `prometheus/client_golang` 暴露 `/metrics` 接口对标，同时可通过 Go 1.21 引入的 `runtime/metrics` 做内存分配探针。
