# Standalone IM Architecture

本文档描述独立 Go + Next.js IM 项目的目标架构。核心目标是让 IM 能在本仓库内完成构建、测试、部署和运维，并通过 connector/provider 机制接入外部平台。

## 1. 总体分层

```text
Next.js Web
  - Workbench
  - Admin Console
  - Diagnostics

Go API
  - Auth / Session
  - Conversation / Message
  - Task / Automation
  - Admin / Config
  - Realtime Gateway

Domain Core
  - Tenant / User / Agent
  - Contact / Conversation / Message
  - OutboundTask / DeliveryReceipt
  - AutomationRun / AuditEvent

Runtime
  - Workers
  - Outbox Relay
  - Connector Runtime
  - Provider Runtime
  - Scheduler

Infrastructure
  - MySQL
  - Redis Streams / PubSub / Cache
  - Object Storage
  - Metrics / Logs / Traces
```

## 2. 消息通道 Connector

消息收发必须通过统一 connector contract：

```text
InboundEvent
OutboundMessage
DeliveryReceipt
ContactIdentity
ConversationBinding
MediaAttachment
```

Connector 职责：

- 解析外部回调并转换为标准 `InboundEvent`。
- 把标准 `OutboundMessage` 转换为外部平台请求。
- 上报送达、失败、撤回、删除、会话变更等 receipt。
- 维护外部身份与本地 contact/conversation 的绑定关系。
- 本仓库内置 `internal.webhook` fake connector，用于不依赖真实平台的本地、CI 和 readiness smoke。

Core 职责：

- 分配会话。
- 存储消息。
- 创建 outbound task。
- 发布 outbox event。
- 推送 realtime event。
- 记录审计。

`contracts/v1/connector-inbound-event.schema.json`、`connector-outbound-message.schema.json` 和 `connector-delivery-receipt.schema.json` 是首批通道中立 connector contracts。企微是 connector，不是 core。后续可加入 Web chat、短信、邮件、内部测试通道和其他 IM 平台。

## 3. RPA 与自动化 Provider

自动化能力必须通过 provider contract：

```text
AutomationCapability
AutomationSession
AutomationTarget
AutomationCommand
AutomationResult
AutomationHealth
```

Provider 职责：

- 声明能力：输入、点击、截屏、上传、下载、音视频、浏览器操作、设备操作。
- 执行命令并返回结构化结果。
- 暴露健康、容量、版本和错误分类。
- 隔离供应商协议、凭证和 runtime。

Core 职责：

- 创建 `AutomationRun`。
- 管理状态机、重试、超时、DLQ 和审计。
- 把自动化结果转换为业务事件或任务终态。

魔云腾/MytRpc 是 provider，不是自动化模型。默认运行面应能用 fake provider 或 HTTP provider 完成验证。

## 4. 写路径

推荐写路径：

1. API 校验输入、认证和权限。
2. API 写入事务性 task/message/receipt。
3. API 写入 outbox。
4. Worker 消费队列并调用 connector/provider。
5. Worker 按幂等键更新终态。
6. Outbox relay 发布 realtime event。
7. Next.js 订阅事件并回源刷新。

要求：

- API 写路径不可直接阻塞在外部平台或设备控制上。
- 每个写任务都有幂等键、状态机、超时、重试和失败原因。
- 外部调用失败必须可补偿、可重放、可观测。

## 5. 实时与投影

Realtime gateway 只发布产品级事件：

- `conversation.created`
- `conversation.updated`
- `message.received`
- `message.sent`
- `message.delivery_updated`
- `task.status`
- `contact.updated`
- `automation.status`

投影表服务于查询和页面体验，不应成为唯一事实来源。Worker 更新投影时要保证：

- 幂等。
- 可重放。
- 可从 canonical event 或主表重建。
- gap 检测后能回源刷新。

## 6. 高可用设计

运行面目标：

- API stateless，可水平扩展。
- Worker 独立部署，可按队列和能力扩缩容。
- Redis Streams 承担缓冲、pending、重试和 backpressure。
- Outbox 确保 DB commit 后事件可发布。
- Cache 失效不影响核心写入。
- 外部 provider 故障不影响 API 接收任务。

可靠性要求：

- 所有外部回调可去重。
- 所有 outbound task 可重试。
- 所有关键状态可观测。
- 所有发布 profile 有回滚路径。
- 关键数据有备份、恢复和一致性检查。

## 7. 安全与多租户

- 所有 API 需要明确租户、用户、角色和权限。
- Connector/provider secret 按租户或企业隔离存储。
- 后台任务必须携带 tenant scope。
- 审计日志记录 actor、action、target、request id 和结果。
- 管理台高风险操作需要额外确认或权限约束。

## 8. 需要从 Core 移出的能力

- 设备协议、屏幕控制、RTC 控制和供应商命令。
- 单一消息平台回调结构。
- 单一平台代理 payload。
- 单一语音转写供应商配置。
- 所有只为过渡期存在的桥接 sidecar。

这些能力可作为 integration 保留，但必须通过 connector/provider contract 接入。
