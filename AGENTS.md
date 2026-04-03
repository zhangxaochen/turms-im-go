# AGENTS.md — Turms IM Go Refactor

> 本文档是 Turms IM Go 重写项目的**主记忆入口**。
> 所有详细设计、分析、规范均分布于 `docs/` 下，由此处索引。

## 项目概览

- **原始项目**: [turms-im/turms](https://github.com/turms-im/turms) — Apache 2.0
- **原始技术栈**: Java 21 + Spring Boot (Reactive) + Netty + MongoDB Sharded + Redis + Protobuf
- **目标**: 用 Go 重写核心服务端 (turms-gateway / turms-service / turms-server-common)
- **原始代码规模**: ~2542 Java 文件, ~690K LOC (含 protobuf 生成代码、客户端)
- **服务端核心代码**: ~300K LOC (turms-gateway + turms-service + turms-server-common)

## 文档索引

| 文档 | 用途 |
|------|------|
| [docs/architecture.md](docs/architecture.md) | 原始架构全景分析 + Go 映射策略 |
| [docs/domain-models.md](docs/domain-models.md) | 数据模型 (MongoDB collections) 完整定义 |
| [docs/business-logic.md](docs/business-logic.md) | 核心业务逻辑详解 (消息 / 群组 / 用户 / 会话 / 会话 / 关系) |
| [docs/infra-components.md](docs/infra-components.md) | 基础设施组件 (集群 / RPC / 插件 / 配置 / ID 生成) |
| [docs/api-protocol.md](docs/api-protocol.md) | 客户端&管理端 API 协议定义 + ResponseStatusCode 体系 |
| [docs/go-design.md](docs/go-design.md) | Go 版总体设计 + 包结构 + 技术选型 |
| [docs/tdd-bdd-strategy.md](docs/tdd-bdd-strategy.md) | TDD + BDD 测试策略详细设计 |
| [docs/review_checklist.md](docs/review_checklist.md) | **Review 固化清单**: 复刻并超越的质量保障准则 |

## 核心约束

1. **读扩散模型**: 消息写一份到 MongoDB，通过 gateway 实时推送给在线用户
2. **无状态服务**: turms-service 完全无状态，turms-gateway 持有 session 状态仅在内存
3. **MongoDB 分片**: 以 `deliveryDate` 做 shard key 支持冷热分离 (tiered storage)
4. **自研 RPC**: 服务端间使用自研二进制 RPC，非 gRPC
5. **Protobuf 编码**: 客户端 ↔ gateway 通信用 protobuf
6. **Snowflake ID**: 全局 ID 生成基于 snowflake 变体
7. **插件系统**: 支持 Java JAR + JavaScript (GraalVM) 双插件
8. **响应式全链路**: 底层全 Reactor/Netty 异步无阻塞

## Go 重写核心原则

1. **保持架构一致**: gateway + service 分离，MongoDB + Redis 存储层不变
2. **Go 惯用法**: interface 替代 Java 继承/注解，struct 替代 class
3. **显式优于隐式**: 不用 DI 框架，手动构造依赖图
4. **性能第一**: goroutine 替代 Reactor，zero-alloc 序列化
5. **可测试性优先**: 所有 service 通过 interface 注入，方便 mock

## [CRITICAL REFRACTOR PROTOCOL]
When porting/refactoring code from Java to Go, you MUST NOT just translate method signatures blindly.
BEFORE committing any ported struct/service, you MUST complete this implicit checklist:
1. **Config Audit**: Which configuration/Properties was the Java version reading? You must read the SAME property in Go using dependency injection. If it does not exist, explicitly mock it or flag it.
2. **Overload Audit**: Java might have multiple overloaded methods with the same name. Go does not have method overloading. You MUST create multiple explicit Go functions (e.g., `NewX`, `NewXWithReason`, `NewXFromError`) to catch EVERY original path. 
   - **Exception**: If a single Go function uses variadic arguments or generic `interface{}` to encapsulate the permutations (e.g. `UpsertSettings`), it is **expected and required** to stack multiple `// @MappedFrom` annotations above that single Go function to track all Java counterparts.
3. **Data Loss Audit**: Ensure all object field assignments via chained Builders in Java (e.g., `.setTimestamp(...)`) are exactly represented in the Go struct literal. Do not drop implicit assignments like system timestamps.
4. **Context Stubbing (TODOs)**: If underlying infrastructure (e.g. cluster RPC, plugins, complex session states) is not yet ported, do NOT just omit the original Java logic implicitly to save time. You MUST add explicit `TODO` comments detailing exactly what Java logic was bypassed and what needs to be wired up later.

## 自动化重构工具 (Automation Tools)

为了追踪映射进度并进行双向挂载，本项目基于 `docs/refactor_progress_report.md` 构建了半自动流水线：

1. **Markdown 链接规范**: 在 `refactor_progress_report.md` 中记录映射关系时，Go 代码的路径和方法必须使用相对路径的 Markdown 链接（例如 `- [x] `javaMethod()` -> [internal/xxx.go:Method()](../internal/xxx.go)`），严禁使用反引号。这允许 GitHub 和 IDE 的直接点击跳转。
2. **自动注入脚本 (`inject_mapped_from.py`)**: 无论何时更新了 `refactor_progress_report.md`，都应该在根目录运行 `python3 inject_mapped_from.py`。该脚本会解析报告中的映射，并自动在对应的 Go 文件函数上方扫描并注入或更新 `// @MappedFrom` 注解，彻底免除手动双写同步的烦恼。
5. **批量审计脚本 (`audit.js`)**: 用于在本地批量验证 `refactor_progress_report.md` 中的所有已打钩项。由于 `gemini` CLI 会遭遇 `429 Too Many Requests` API 频率限制，该脚本通过延时重试策略进行单次 Query 的并行控制，检测出的问题会自动归档到 `pending_bugs.md`。

## 代码映射排坑经验 (Gotchas & Insights)

- **避免单体 Stub 文件**: 在映射原本庞大且分散的 Repository 或 Controller 时，**绝不能**使用一个集中的巨石文件（例如 `group_repositories_stubs.go`）来堆砌所有接口存根。必须精准将映射还原到其属于的特定 Domain Repository 文件（例如 `group_blocked_user_repository.go`），以保证依赖层次与 Java 原始设计解耦的初衷完全一致。
- **批量查询限流**: 在编写或运行利用大型语言模型代理（如使用 gemini CLI）进行自动化架构校验的脚本时，切记限制并发请求并结合失败/限流（429 报错）时的长时间回退重试（Exponential Backoff）。建议不要把整个代码库压缩到一个 Prompt 中，应当以 Class 为单位发起单次 Query (`-p` 参数)，降低单请求上下文并保证稳定性。
- **命令行工具代理读取**: 使用 `gemini` CLI 进行自动化文件比对或审核时，**严禁通过 `cat` 或 `stdin` 将文件内容拼接到 Prompt 中**（极易触发 Node `execSync` 缓存限制或 Shell 字符数限制）。正确做法是直接在 Prompt 字符串中提供文件的绝对路径（如 `gemini -p "Read /path/a and /path/b..."`），使 AI Agent 借助内置的 `view_file` tool 自行读取文件，极大提升脚本稳定性和安全裕度。
