# 核心业务逻辑解析 (Business Logic)

Turms 的业务逻辑高度集中在 `turms-service` 模块下。我们按照四大核心领域进行梳理，为 Go 迁移提供实现指引。

## 1. 消息模块 (`MessageService`)
处理即时通讯的最核心职责。

### 核心流转规则
1. **鉴权与校验**: 发送消息前验证发送者状态（是否被封禁，是否在线），校验 `targetId`（单聊是否好友/拉黑，群聊是否群成员）。
2. **ID 统一分配**: 
   - 调用 Snowflake Node 生成全局消息 ID (`messageId`)。
   - 对严格序要求的会话，调用 Redis 分配递增的 `SequenceID`。
3. **存储持久化**: 异步写入 MongoDB 的 `message` collection。
4. **实时推送 (Routing/Pushing)**: 
   - 计算消息接收人列表（单聊为对方，群聊为查询出的全体有效群成员）。
   - 只给通过 Redis 中查得**当前在线**的 Session 推送。
   - 通过 RPC 调用目标 `turms-gateway` 节点把 Payload 下发到连接上。
5. **缓存击穿防护**: `MessageService` 大量使用了本地缓存 Caffeine。Go 迁移时可使用 `dgraph-io/ristretto` 或 `puzpuzpuz/xsync` + 局部 `lru` 来替代。

## 2. 会话模块 (`SessionService`)
位于 Gateway 和 Service。

### Gateway 层的 `SessionService` (状态持有者)
- **连接打点**: TCP/WebSocket 建立连接并鉴权后，在内存 (ConcurrentHashMap) 中注册 `Session` 对象，维护心跳时间。
- **冲突策略**: 遇到多端同时登录相同设备类型，触发下线旧设备 (conflict resolution)——并向旧设备推送 `SESSION_SIMULTANEOUS_CONFLICTS_OFFLINE` 断线帧。

### Service 层的 `SessionService` (全局状态管理者)
- 依赖 Redis 存储全局用户的在线状态、所在的 gateway 节点 IP 和端口。
- 供其他 Service 模块查询“某用户当前连在哪个 Gateway”以进行 RPC 投递。

## 3. 群组管理 (`GroupService`)
- 包含子领域：`GroupMemberService` (成员管理), `GroupJoinRequestService` (加群申请), `GroupBlocklistService`。
- **架构特点**: 大量依赖 MongoDB 的原子操作 (如 `$addToSet`, `$pull`, `$set`) 以避免并发更新导致的数据不一致。Go 中必须保持这些原子 BSON 更新操作。

## 4. 离线消息拉取 (Pulling)
- 客户端上线并不会收到堆积的未读离线消息全量推送，而是主动带着上次同步的时间戳 (或者 lastMessageId) 发起 Pull 请求。
- `MessageService` 结合 `deliveryDate` 索引去 MongoDB 捞取该时间节点后的增量消息。
