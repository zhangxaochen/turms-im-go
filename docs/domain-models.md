# Turms 数据模型定义 (Domain Models)

Turms 的存储层主要依赖 MongoDB，且大量利用了 MongoDB 的分片 (Sharding) 和复合索引 (Compound Indexes)。在转译为 Go 时，需要通过 `go.mongodb.org/mongo-driver/bson` tags 进行精确映射，并在服务启动时通过代码确保相应的 Index 被正确创建。

## 1. 核心 Collection 模型

### 1.1 消息模型 (`Message`)
代表私聊、群聊和系统消息。
- **Collection Name**: `message`
- **Shard Key**: `deliveryDate` (为了支持基于时间的冷热分离 Tiered Storage)
- **ID**: `Long` (非主要查询条件，主用组合索引)
- **关键索引设计**:
  - `(deliveryDate, targetId)`: 私聊和群聊的通用组合索引。
  - `(deliveryDate, conversationId)`: （只有存在 conversationId 时索引），用于快速检索特定会话的消息。

**Go Struct 映射示例**:
```go
type Message struct {
	ID               int64     `bson:"_id"`
	ConversationID   []byte    `bson:"cid,omitempty"`  // cid
	IsGroupMessage   *bool     `bson:"gm,omitempty"`   // gm
	IsSystemMessage  *bool     `bson:"sm,omitempty"`   // sm
	DeliveryDate     time.Time `bson:"dyd"`            // dyd
	ModificationDate time.Time `bson:"md,omitempty"`   // md
	DeletionDate     *time.Time`bson:"dd,omitempty"`   // dd
	RecallDate       *time.Time`bson:"rd,omitempty"`   // rd
	Text             string    `bson:"txt,omitempty"`  // txt
	SenderID         int64     `bson:"sid"`            // sid
	SenderIP         *int32    `bson:"sip,omitempty"`  // sip
	SenderIPv6       []byte    `bson:"sip6,omitempty"` // sip6
	TargetID         int64     `bson:"tid"`            // tid
	Records          [][]byte  `bson:"rec,omitempty"`  // rec
	BurnAfter        *int32    `bson:"bf,omitempty"`   // bf
	ReferenceID      *int64    `bson:"rid,omitempty"`  // rid
	SequenceID       *int32    `bson:"sqid,omitempty"` // sqid
	PreMessageID     *int64    `bson:"pmid,omitempty"` // pmid
}
```

### 1.2 群组模型 (`Group`)
代表一个聊天群组的元数据。
- **Collection Name**: `group`
- **Shard Key**: `_id` (通常是 Hash Sharding)
- **关键字段**: `creatorId`, `ownerId`, `name`, `intro`, `minimumScore`, `isActive`
- **自定义属性**: `userDefinedAttributes` (在 Go 中映射为 `map[string]interface{}`)
- 有过期和封禁机制 (依赖 `deletionDate` 或 `muteEndDate`)。

### 1.3 用户与关系模型
- **`User`**: 存储用户的基本信息 (密码哈希、姓名、注册日期等)。
- **`UserRelationship`**: 表示双向好友关系、黑名单等。关联 `UserRelationshipGroup`。

---

## 2. 字段生命周期与可扩展性
- **软删除设计**: 如 `Message`, `Group` 都带有 `deletionDate`。只要字段存在，该实体通常视为已被逻辑删除 (Expirable)。
- **Customizable**: Turms 极度强调定制能力，像 `Group` 这种均实现了 `Customizable` 接口，这在 Go 中表现为嵌入式结构体或统一挂载 `CustomAttributes map[string]any`。

## 3. ID 与序列生成
- **Message ID**: Turms 不依赖 MongoDB ObjectId，而是通过 Redis 配合本地时钟生成的 Snowflake ID，确保全局唯一及时间递增特性。
- **Sequence ID (`sqid`)**: 用于保证消息的严格递增顺序，依赖 Redis Incr 实现。Go 中必须复刻此发号器逻辑。
