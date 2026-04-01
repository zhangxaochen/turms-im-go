# TDD 与 BDD 混合测试策略 

作为基础中间件重写项目，容错率极低。必须保障对历史 API 兼容性的 100% 验证。为此，设计如下维度的自动化测试约束：

## 1. 单元测试 TDD (Unit & Behavior)
业务逻辑主要集中在 `internal/domain`。对所有 `Service` 必须进行基于 Behavior 的行为验证。
- **框架支持**: 使用标准 `testing` + `stretchr/testify` (包含 `assert` 与对服务的依赖 `mock`)。
- **Mock 策略**: 
  - 所有 DB操作、RPC操作必须通过接口 `Interface` 访问以方便 `testify/mock`。
  - 核心逻辑 (例如《群发消息发送失败的回滚》、《离线拉取游标更新机制》) 使用 Table-Driven Tests (表格驱动测试)。

示例：
```go
func TestMessageService_SendMessage(t *testing.T) {
    tests := []struct{
        name        string
        setupMock   func(*MockMongoDao, *MockRedisDao)
        request     *pb.CreateMessageRequest
        wantErr     bool
        errCode     int
    } {
        // ... 详尽罗列场景
    }
}
```

## 2. 真实集成测试 BDD (Integrations via Testcontainers)
模拟服务运行时与底层系统的交互，避免“由于 Mock 导致的问题假象”。
- **基础设施引入**: 引入 `testcontainers/testcontainers-go`。在测试运行之前：
  1. 调用 Docker API 即时拉起一份干净隔离的 `MongoDB 6.x` 和 `Redis 7.x` 容器。
  2. 初始化必需的 TTL 索引和 Compound Index。
- **测试范围**: 只涉及 `storage/` 下的方法以及复杂查询逻辑，确保语句与聚合 (Pipeline) 均正确无误。

## 3. 端到端 E2E (Client Simulation Sandbox)
验证最终 RPC/网络协议边界层的完美映射。
- **黑盒跑批**:
  启动 `turms-service` 与 `turms-gateway`。
  利用现有 Java 测试工具/脚本（或者我们写的 Go `net/tcp` 裸协议小脚本）建立多条连接，伪造客户端登录和消息交互。
  确认状态码：
    1. Java 客户端断联重连逻辑正常触发。
    2. 多端登录互踢收到的 ResponseStatusCode 以及 Notification Packet 正确。

## 4. 压力保障边界 (Benchmarks)
- 在 CI 流程中加入 `go test -bench=. -benchmem`，专门针对高频执行的发包和拆包编解码代码进行性能约束。
- 只有内存分配次数 (allocs) 未增长才能被 review 并 merge，否则阻止退化。
