# 移动端 App 蓝图（iOS / Android）

> 状态：设计基线（尚无代码） · 隶属计划：[移动端 App P0](2026-07-11-mobile-app.md) · 创建：2026-07-11 · 更新：2026-07-11

本文档是移动端 App 的架构与功能规范，供后续 AI 代理（Claude/Codex）据此开发。它描述目标
设计，不是当前事实；后端与 Web 的当前事实以[架构总览](../../architecture/overview.md)、
[功能地图](../../product/feature-map.md)和代码为准，冲突时先核对代码再修正本文档。
栈选型与边界决策的原因见 [ADR-0006](../../decisions/0006-mobile-app-thin-client.md)。

## 产品定位

移动端是「商机处理端」，不是 Web 的全量移植。核心价值是 Web 端缺失的一环：**推送触达 +
随手处理**（收到商机 → 审核 → 回复 → 跟进，全程在手机完成）。规则管理、模板编辑、企业微信
绑定等重配置场景留在 Web。

## 技术栈与仓库位置

| 决策项 | 选择 | 说明 |
| --- | --- | --- |
| 代码位置 | 仓库根 `mobile/`，与 `frontend/`、`backend/` 并列 | 单仓多端，契约同步在同一变更内完成 |
| 框架 | React Native + Expo（managed）+ TypeScript strict | 单代码库出 iOS/Android，见 ADR-0006 |
| 路由 | expo-router | 文件式路由，支持推送深链 |
| 服务端状态 | TanStack Query | 列表轮询、详情缓存、失败重试统一处理 |
| 本地状态 | 组件 state / React Context | 不建全局 store；禁止复刻 Web `AppStoreProvider` 的演示态混合模式 |
| 凭据存储 | expo-secure-store | JWT 只存 SecureStore，不进 AsyncStorage/日志 |
| 推送 | expo-notifications + Expo Push Service | 退出路径为直连 APNs/FCM，见 ADR-0006 |

## 与后端的契约

- App 直连 v1 REST API + JWT Bearer，不建 BFF。
- `mobile/lib/api/` 是唯一 HTTP 边界，`mobile/lib/types.ts` 镜像后端 DTO，与
  `frontend/lib/api.ts`、`frontend/lib/types.ts` 同模式。
- 契约同步链：后端 Pydantic DTO → `frontend/lib/types.ts` → `mobile/lib/types.ts` → UI。
  后端字段变化必须三端连贯更新（[功能地图](../../product/feature-map.md)同步清单的扩展）。
- API 具体查询参数、错误码以 `backend/app/api/v1/routes/` 与 DTO 为准，本文档只锚定端点。

## 目标目录结构

```
mobile/
  app/                 # expo-router 路由（登录、收件箱、详情、设置）
  features/
    inbox/             # 商机收件箱：列表、筛选、刷新
    opportunity/       # 详情、消息历史、Agent 发现、回复、状态流转
    settings/          # 订阅只读、连接引导、关于
  lib/
    api/               # 唯一 HTTP 边界（fetch 封装 + 各资源 client）
    types.ts           # 后端 DTO 镜像
    auth/              # OAuth、token 存取、会话恢复
    push/              # 推送注册、深链解析
  components/          # 跨 feature 复用 UI
```

## 功能模块

### P0 — 最小可用闭环

除推送外全部依赖已存在的后端能力；成熟度结论引自[功能地图](../../product/feature-map.md)。

| # | 模块 | 依赖 API | 后端现状 | 验收要点 |
| --- | --- | --- | --- | --- |
| 1 | 登录 | `GET /auth/oauth/{provider}/authorize`、`POST /auth/oauth/{provider}/native`（待新增）、`GET /auth/me` | OAuth/JWT 真实；原生 id_token 端点缺失 | Google/Apple 原生登录；token 存 SecureStore；冷启动经 `/auth/me` 恢复会话 |
| 2 | 商机收件箱 | `GET /opportunities` | 已实现，owner 隔离 | 列表、状态/渠道筛选、分页、下拉刷新；P0 先轮询，推送落地后改为推送触发刷新 |
| 3 | 商机详情 | `GET /opportunities/{id}`、`GET /messages` | API 已实现；Web 未消费这两个端点，不得照抄 Web 本地 store 实现 | 独立请求详情与消息历史；展示检测结果与 Agent 发现（链接核验结论、联系方式、紧急标记） |
| 4 | 回复 | `POST /opportunities/{id}/manual-reply`、`POST /opportunities/{id}/ai-draft`、`GET /templates` | 后端真实发送/生成/落库 | 手动回复真实发送；AI 草稿可编辑后发送；模板只读选用；发送失败可重试且不得伪造已回复状态 |
| 5 | 状态流转 | `PATCH /opportunities/{id}/status`、`POST /opportunities/{id}/claim` | 后端已实现（领域状态机约束） | 认领/跟进/关闭走后端状态机；非法迁移展示后端错误而非本地吞掉 |
| 6 | 推送通知 | 设备注册 API + Celery 发送任务（均待新增） | 后端缺失，唯一新增面 | 新商机、重大商机提醒、AI 自动回复结果、额度耗尽四类事件；点通知深链进详情；`PUSH_ENABLED=false` 时全链路安全关闭 |

### P1 — 差异化能力

| # | 模块 | 依赖 | 说明 |
| --- | --- | --- | --- |
| 7 | Telegram 连接引导 | 统一连接 API（`features/telegram-native-connections` 分支合入后可用） | app 深链拉起 Telegram 完成群授权后回跳；移动端体验优于 Web |
| 8 | Agent 动作审批 | 商机详情中的 Agent 投影；批准后的执行用例后端暂无 | 只读建议列表 + 批准/驳回意向落库；未经独立审批用例不得触发任何外部发送 |
| 9 | 订阅与用量 | `GET /subscriptions/plans`、`GET /subscriptions/me` | 只读展示套餐、AI 额度、TG 限额；升级引导跳 Web（规避 IAP 抽成与审核复杂度） |
| 10 | 今日概览 | `GET /stats/summary` | API 已实现且当前无消费方；做商机数/回复率小面板 |

### P2 — 后置

- 通知偏好与每日摘要：需要新增用户通知偏好 API（Web 端此功能也仅是演示态）。
- 工作时间快捷开关：`/configs/work-mode` 已有，低频操作。
- 离线只读缓存：TanStack Query persist，有真实需求再做。

### 明确不移植

- 好友申请执行按钮（Web 演示态，timer 模拟，无后端）。
- Web 的通知偏好开关（仅页面 state）。
- 规则管理、模板编辑、企业微信绑定（留在 Web）。

## 后端新增面

移动端要求的后端改动只有两块，均需遵守[开发规范](../../development/standards.md)与
[安全基线](../../quality/security.md)：

1. **移动登录端点（P0）**：`POST /auth/oauth/{provider}/native` 校验原生登录取得的
   id_token（复用 `app/core/security.py` 的 JWKS 验签），签发与 Web 相同的 JWT。现有
   callback 面向 Web 重定向，不适配 app。
2. **推送通道（P0）**：DeviceToken 模型 + Alembic 迁移、注册/注销 API（按用户隔离）、
   Celery 发送任务挂接三个既有事件点（摄取用例的审核者通知、重大商机 attention 投影、
   AI 自动回复结果），新增 `PUSH_ENABLED` 环境变量作为与 `IM_SEND_ENABLED` 同级的安全阀。
3. **通知偏好 API（P2）**：与 app 通知设置页同期设计。

## 安全与不变量

承接[安全基线](../../quality/security.md)与架构总览的主要不变量，移动端追加：

- JWT 只存 SecureStore；日志与崩溃上报不得包含 token、消息内容、联系方式。
- 推送 payload 最小化：只含事件类型 + 资源 ID，打开 app 后凭 JWT 拉取详情。
- owner 隔离由后端保证；app 切换账号时必须清空本地缓存与查询缓存。
- 深链参数先校验（资源 ID 格式、当前登录态）再路由。
- 发送失败不得在 UI 伪造已回复；AI 操作入口展示剩余额度，额度耗尽（后端 fail-closed）
  给出明确提示而非静默失败。

## 验证入口（随 P0 落地）

- `make mobile-check`：`tsc --noEmit` + ESLint + 单元测试，与 `frontend-check` 并列进 CI。
- 每个切片的完成定义：接真实 API、处理 loading/error/空态、通过 `mobile-check`。
- 进度、发现与决策写回[执行计划](2026-07-11-mobile-app.md)，不留在会话里。
