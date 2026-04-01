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
