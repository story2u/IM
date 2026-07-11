# 移动端 App P0（iOS/Android 商机处理端）

> 状态：active · Owner：bruce / AI 代理 · 创建：2026-07-11 · 更新：2026-07-11

## 目标与用户价值

商机 owner 在手机上实时收到商机推送，并完成「审核 → 回复 → 跟进」闭环，覆盖 Web 端因
轮询与桌面在场限制而错过的时段。交付物是 `mobile/` 下可运行的 iOS/Android app（P0 模块）
与配套的后端推送通道。设计规范见[移动端 App 蓝图](2026-07-11-mobile-app-blueprint.md)，
栈决策见 [ADR-0006](../../decisions/0006-mobile-app-thin-client.md)。

## 非目标

- 规则管理、模板编辑、企业微信绑定（留在 Web）。
- app 内订阅付费/升级（IAP），P1 仅只读展示。
- 好友申请等 Agent 外部动作的实际执行。
- 离线优先与富媒体推送。

## 背景与当前行为

- 仓库当前没有 `mobile/` 目录；后端无推送通道，摄取用例的「通知审核者」只到队列接口。
- P0 所依赖的商机、消息、回复、状态、模板 API 均已实现（见[功能地图](../../product/feature-map.md)），
  但 Web 端多处未消费（详情/消息历史），存在演示态本地 store —— 移动端不得照抄。
- OAuth 现状是面向 Web 的 authorize/callback 重定向（`backend/app/api/v1/routes/auth.py`），
  callback 跳转 `frontend_base_url/login#token=...`，app 无法直接复用。

## 验收标准

- [ ] Google/Apple 原生登录换取后端 JWT，冷启动自动恢复会话，登出清空 SecureStore 与缓存。
- [ ] 收件箱展示当前用户商机，支持状态/渠道筛选与分页刷新；无跨用户数据。
- [ ] 详情页独立请求 `GET /opportunities/{id}` 与消息历史，展示检测结果与 Agent 发现。
- [ ] 手动回复经后端真实发送并展示 outgoing 落库结果；失败可重试且不伪造已回复状态。
- [ ] AI 草稿生成可编辑后发送；额度耗尽时展示明确提示（后端 fail-closed）。
- [ ] 状态流转与认领走后端状态机，非法迁移展示后端错误。
- [ ] 新商机、重大商机提醒、AI 自动回复结果、额度耗尽四类事件产生推送；点通知深链进详情。
- [ ] `PUSH_ENABLED=false` 时不注册、不发送、不报错；推送 payload 不含消息内容。
- [ ] Web 演示态能力（好友申请执行、通知偏好开关）未被移植。
- [ ] `make mobile-check` 与后端相关测试在 CI 通过。

## 影响面与风险

- 新增 `mobile/` 目录（Expo app），根 `Makefile` 与 CI 增 `mobile-check`。
- 后端：新增原生登录端点、DeviceToken 模型 + Alembic 迁移、注册/注销 API、Celery 推送
  任务、`PUSH_ENABLED` 配置；涉及认证与外发通道，按高风险变更走[安全基线](../../quality/security.md)。
- frontend 无改动；domain 层预计无改动（推送编排放 application/worker）。
- 风险：Expo Push 第三方投递可靠性（退出路径见 ADR-0006）；App Store 审核对登录/推送的
  合规要求；契约三端同步成本上升。

## 实施步骤

- [ ] 1. `mobile/` Expo 脚手架 + TypeScript strict + `make mobile-check` + CI 接入。
- [ ] 2. 后端 `POST /auth/oauth/{provider}/native`（id_token 验签复用既有 JWKS 逻辑）+ 测试。
- [ ] 3. app 登录/会话恢复/登出（SecureStore + `/auth/me`）。
- [ ] 4. 收件箱列表（轮询版）+ 筛选分页。
- [ ] 5. 详情页：详情 + 消息历史 + Agent 发现展示。
- [ ] 6. 回复：手动回复、AI 草稿、模板只读选用。
- [ ] 7. 状态流转与认领。
- [ ] 8. 后端推送通道：DeviceToken 迁移 + 注册 API + Celery 任务挂接三个事件点 + `PUSH_ENABLED`。
- [ ] 9. app 推送注册/接收/深链，收件箱改推送触发刷新。
- [ ] 10. 文档回写：功能地图增移动端条目、运维文档增环境变量、归档本计划。

每步为可独立验证的纵向切片；步骤 2、8 涉及认证与外发，需相称的后端测试。

## 进度日志

- 2026-07-11：创建蓝图、ADR-0006（proposed）与本计划；等待评审后从步骤 1 开始。

## 发现日志

- 2026-07-11：`auth.py` callback 面向 Web 重定向（`frontend_login_redirect`），确认移动端
  需要独立的原生 id_token 校验端点，已写入蓝图「后端新增面」。
- 2026-07-11：`GET /stats/summary`、`GET /messages`、`GET /opportunities/{id}` 已实现但 Web
  未消费，移动端可直接使用，无需等待 Web 改造。

## 决策日志

- 2026-07-11：栈与边界（RN/Expo 瘦客户端、无 BFF、Expo Push 起步、原生登录端点）见
  ADR-0006，状态 proposed；分支合并采用后改 accepted。
- 2026-07-11：ADR 编号取 0006——0005 已被 `features/telegram-native-connections` 分支的
  统一连接模型 ADR 占用，避免合并冲突。

## 验证记录

| 命令/场景 | 结果 | 证据或备注 |
| --- | --- | --- |
| `make harness-check` | 通过（2026-07-11） | 30 个 Markdown 链接完整、无孤儿文档，65 个后端 Python 文件边界检查通过 |

## 回滚与恢复

- 当前阶段为纯文档，revert 提交即可。
- 后续 `mobile/` 为独立目录，不影响现有部署单元；后端推送通道受 `PUSH_ENABLED` 安全阀
  控制，可配置级关闭；DeviceToken 迁移提供 downgrade。

## 结果与剩余风险

进行中；完成时补实际交付、偏差与后续链接。
