# Turms Refactor Review Checklist (固化)

为了确保 Turms Go 重写项目在逻辑、性能和稳定性上“复刻并超越”Java 原生实现（turms-orig），在 Review 所有包含 `[x]` 的已完成项时，必须严格执行以下检查清单。

## 1. 核心业务与逻辑对齐 (Logic Parity)
- [ ] **完整性**: 是否涵盖了 Java 原生 Service 中的所有公共方法和核心私有逻辑？
- [ ] **防御性编程**: 
  - 是否对非空指针（Pointer Check）进行了校验？
  - 是否对输入的 Protobuf 请求进行了验证（使用 `internal/infra/validator`）？
- [ ] **鉴权与状态控制**: 
  - 是否复刻了针对“用户是否被封禁”、“是否为好友”、“权限等级”等的业务检查？
- [ ] **多端一致性**: 如果逻辑涉及多端登录冲突，处理策略是否与 Java 版一致？

## 2. 数据存储与事务 (Persistence & Transactions)
- [ ] **原子更新**: 
  - 是否使用了 MongoDB 的原子操作符（`$set`, `$inc`, `$addToSet`, `$pull`）？
  - 严禁在没有分布式锁的情况下进行“先读后写”的非原子操作。
- [ ] **事务支持**: 
  - 在涉及多集合更新（如关系 + 版本号）或强一致性场景下，是否正确传递并使用了 `mongo.Session`？
  - `turmsmongo.ExecuteWithSession` 是否能正确处理回滚？
- [ ] **索引效率**: 新增查询逻辑是否利用了现有复合索引？是否避免了全表扫描？

## 3. 并发、内存与性能 (Concurrency & Memory)
- [ ] **锁竞争**: 
  - 是否避免了对共享全局 Map 的单一锁竞争？
  - 高频访问的 Session/Relationship 字典是否使用了分片锁（Sharded Map）？
- [ ] **内存分配**: 
  - 核心读写路径（Hot Path）是否使用了 `sync.Pool` 减少 GC 压力？
  - 是否避免了在大循环中进行不必要的 slice/map 创建？
- [ ] **Context 处理**: 是否正确收敛并传递了 `context.Context` 以支持超时取消和链路追逐？

## 4. 缓存管理 (Cache Management)
- [ ] **缓存失效**: 修改数据（Upsert/Delete）后，是否**显式**调用了对应的缓存失效逻辑（如 `invalidMemberCache`）？
- [ ] **击穿与一致性**: 
  - 本地缓存（Local Cache）是否配置了合理的过期时间？
  - 是否实现了防止缓存击穿（Cache Breakdown）的加载保护？

## 5. 版本与同步逻辑 (Versioning)
- [ ] **Version Sync**: 修改操作完成后，是否按需更新了 `UserVersion` 或 `GroupVersion`？
- [ ] **异步更新**: 为了响应速度，版本更新是否在异步 Goroutine 中执行（除非业务要求强一致同步）？

## 6. 协议与错误码 (Protocol & Codes)
- [ ] **错误码对齐**: 抛出的异常是否使用了标准 `pkg/codes` 并在 Controller 中映射为正确的 `ResponseStatusCode`？
- [ ] **数据还原**: 返回给客户端的数据结构是否经过了完整映射，确保字段名和类型在 Protobuf 层面 100% 兼容？

---

## 执行建议
1. **AI Reviewer**: 在每次标记 `[x]` 前，自动运行此清单中的检查项并生成简要报告。
2. **劣化判定**: 若新实现的代码量大幅少于 Java 版，需额外警剔是否遗漏了边界条件（Edge Cases）。
3. **超越判定**: 若新实现利用了 Go 的非阻塞特性（Channels/Select）提升了吞吐量，需在 Walkthrough 中记录性能提升点。
