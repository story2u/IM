# ADR-0006：移动端作为现有 API 的瘦客户端（React Native + Expo）

> 状态：proposed · 日期：2026-07-11

## 背景

商机时效性强，Web 端依赖 30 秒轮询和用户守在桌面前，非工作时间与移动场景缺少触达和处理
能力。后端 v1 REST API 已按 owner 隔离，覆盖商机列表、详情、消息历史、回复、状态流转全流程；
团队现有技术栈是 TypeScript/React（Next.js 前端）。需要同时交付 iOS 与 Android。

## 决策

- 新建 `mobile/` 目录，使用 React Native + Expo（managed workflow）+ TypeScript strict 单代码库
  交付 iOS/Android。
- App 直连现有 FastAPI v1 REST API + JWT，不引入 BFF/GraphQL 聚合层。
- HTTP 访问边界唯一（`mobile/lib/api/`），DTO 类型镜像与 Web `frontend/lib` 同模式；Web 中的
  演示态能力（timer/mock）一律不移植。
- 推送是唯一新增后端能力：DeviceToken 模型 + 注册 API + Celery 发送任务 + `PUSH_ENABLED`
  安全阀；投递起步经 Expo Push Service，payload 只携带类型与 ID，不携带消息内容。
- 移动端 OAuth 采用原生登录 + 后端 id_token 校验端点（复用既有 JWKS 验签），不复用面向
  Web 的 callback 重定向。

## 备选方案

- 原生双端（Swift/Kotlin）：性能与平台一致性最好；两份实现、与团队 TS 栈不符，P0 没有重
  交互场景。未选择。
- Flutter：单代码库；Dart 与现有 TS 契约镜像割裂，类型无法复用。未选择。
- PWA：零上架成本；iOS 推送可靠性与留存差，而推送是本产品移动端的核心价值。未选择。
- BFF/GraphQL：聚合灵活；当前只有一个移动客户端，无聚合需求。未选择。

## 后果

- 契约同步链从两端扩为三端：后端 DTO → Web types → mobile types，功能地图同步清单需随
  移动端落地扩展。
- 依赖 Expo Push 第三方投递；DeviceToken 结构保持 provider 无关，退出路径是直连 APNs/FCM。
- App Store 上架要求 Apple 登录（后端 provider 已具备）；app 内订阅升级受 IAP 规则约束，
  P1 只做只读展示并引导去 Web。
- 推送不含消息明文，通知内容泄露面最小化，但用户需点开 app 才能看到详情。

## 验证与复审

P0 验收（见[移动端 App P0 计划](../plans/active/2026-07-11-mobile-app.md)）端到端跑通即视为
决策有效。复审信号：Expo Push 投递失败率或延迟不可接受、需要富媒体推送时复审直连
APNs/FCM；出现第二个客户端形态或聚合需求时复审 BFF。
