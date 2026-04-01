# 客户端/服务端 API 通信协议 (Protobuf & Codes)

Turms 的协议设计极其严谨和详尽，包含一个庞大的状态码库和数以百计的 Protobuf Message 类型。这部分在迁移到 Go 时基本上可以依靠代码生成来完成，**不需要甚至不可以手工改写**。

## 1. Protobuf 传输载荷 (Payload)
- Turms 使用了单一的巨大请求封包模型 `TurmsRequest`，里面包含了数十种可能的 `OneOf` 内部请求（如 `CreateMessageRequest`, `CreateGroupRequest` 等）。
- 这种设计允许客户端与所有的网关通信只需维护一套 Socket 或者 WebSocket 序列化即可。

### Go 生成策略
```bash
protoc --go_out=./pb --go_opt=paths=source_relative \
       --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
       turms-plugin-livekit/...  turms-client/...
```
**注意**: 生成的 Go struct 指针类型检查一定要注意 `nil` 判断以免 panic（相较于 Java Protobuf 的 Optional）。

## 2. 状态码字典 (`ResponseStatusCode.java`)
高达将近 500 个细分业务错误码 (Business Codes)。每个错误码同时映射至一个标准 HTTP 状态码。

### 原设计特点
```java
// User - Session
SESSION_SIMULTANEOUS_CONFLICTS_DECLINE(2100, "A different device has logged into your account", 409),
SESSION_SIMULTANEOUS_CONFLICTS_OFFLINE(2102, "A different device has logged into your account", 409),

// Message - Send
NOT_FRIEND_TO_SEND_PRIVATE_MESSAGE(5001, "Only the friends or not blocked users can send messages", 403),
```

### Go 处理方式 (代码隔离与错误体系)
在 Go 中，推荐用一个专门的 `errors` 或 `codes` package 定义常量，同时封存一个自定义的 `TurmsError`：

```go
package codes

const (
    InvalidRequest                    = 1100
    SessionSimultaneousConflictsOffline= 2102
    NotFriendToSendPrivateMessage     = 5001
)

type CodedError struct {
    BusinessCode int
    HTTPStatus   int
    Reason       string
}

func (e *CodedError) Error() string {
    return e.Reason
}

// 封装一个预定义的 map 将 Business Code 映射出标准的 CodedError 返回。
func GenerateErr(code int) *CodedError { ... }
```
这样不仅能利用 Go 的原生 `error` 处理惯例，还能完美地对应并投递回 Protobuf `Response` 帧中去，保持对旧客户端 100% 兼容。
