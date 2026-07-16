# 工作机会发现垂直模式

> 状态：active · Owner：Codex · 创建：2026-07-16 · 更新：2026-07-16

## 目标与用户价值

在现有消息摄取、用量账本和受限 pi Agent Runtime 上增加工作机会发现模式。系统从用户已授权
的 Telegram/企业微信来源中识别真实招聘信息，保存可追溯职位字段，并按用户主动声明的职业
偏好做确定性匹配，在 Web、iOS、Android 提供一致的列表和详情体验。

## 非目标

- 不自动投递、联系招聘方、接受 Offer 或执行其他外部动作。
- 不为雇主筛选候选人，不收集或推断受保护属性。
- 不承诺未经真实生产样本验证的识别准确率。
- 不建立第二套 Agent Runtime、用量账本或链接检查器。

## 背景与当前行为

基线为 `release/v2.0.0@e65c9ca`。`Message` 已保存 owner、来源、作者、平台时间和外部消息
标识，`Opportunity` 仅表达通用商机。摄取入口位于
`backend/app/application/use_cases/ingest_message.py`，pi 后处理位于 `analyze_message.py` 和
`backend/pi-agent-runtime/`。三端已有商机列表、详情、认证和统一 API 边界。

## 验收标准

- [ ] 来源职能画像按 owner 隔离、可缓存、可人工覆盖。
- [ ] 预筛、招聘分类和结构化提取不阻塞普通消息摄取，且重复消息不重复计费。
- [ ] `job_post`/`job_repost` 才创建职位，缺失字段保持 null/unknown，重要字段保留证据。
- [ ] 旧商机默认 `business`；职位详情、档案、匹配、反馈均持久化且迁移可回滚。
- [ ] 匹配分由确定性规则计算，年龄/性别等限制只产生合规提示。
- [ ] Web/iOS/Android 提供工作机会列表、详情和档案入口。
- [ ] 虚构 eval 数据和脚本可运行，结果不冒充真实生产准确率。
- [ ] `make check` 及平台检查通过，功能地图和架构文档与代码一致。

## 影响面与风险

- domain：新增职位、来源、分类和匹配枚举及纯策略。
- infrastructure：新增 SQLModel、Alembic、repositories；扩展受限 Node runner schema。
- application/worker：增加异步职位发现用例，复用当前 message Agent 额度和任务队列。
- API：增加 `/jobs`、`/job-search-profiles` 和来源画像端点，所有资源强制 owner 隔离。
- clients：三端 DTO、导航、列表、详情和档案编辑。
- 安全：就业高影响边界、证据约束、原始消息最小展示、禁止受保护属性参与评分。

## 实施步骤

- [ ] 领域枚举、模型、迁移和迁移测试。
- [ ] 来源画像、规则预筛、Agent 分类提取和异步任务。
- [ ] 去重、确定性匹配、档案与 repositories。
- [ ] owner 隔离 API 和 DTO。
- [ ] Web、iOS、Android 客户端。
- [ ] eval、测试、文档和完整验证。

## 进度日志

- 2026-07-16：完成 release 审计并创建功能分支；下一步建立领域模型和迁移。

## 发现日志

- 当前消息已有真实 `sent_at`，职位 `posted_at` 可直接使用，不需要模型猜测。
- 当前 Telegram/企微统一消息只有群名，没有稳定群描述字段；画像模型需允许 description 为空。
- 推送通道仍未交付，本次只保存通知偏好和匹配，不声称已推送职位。

## 决策日志

- 2026-07-16：复用现有 message 级 pi Agent 调用和 UsageLedger；避免一次消息重复模型计费。
- 2026-07-16：匹配是纯领域服务；模型只解析职位/偏好并生成解释，不能写分数或决定资格。
- 2026-07-16：第一版语义去重使用可替换的相似度端口，确定性指纹为生产主路径；没有 embedding
  provider 时不静默伪造语义聚类。

## 验证记录

| 命令/场景 | 结果 | 证据或备注 |
| --- | --- | --- |
| `make backend-check` | 待运行 | |
| `make frontend-check` | 待运行 | |
| `make ios-check` | 待运行 | |
| `make android-check` | 待运行 | |
| `make check` | 待运行 | |

## 回滚与恢复

先关闭职位发现任务入口，再部署前一版本；执行新迁移 downgrade 会删除职位投影、档案、匹配和
画像表，并移除 `opportunities.opportunity_type`。原始 `messages`、通用商机和 IM 连接不受影响。
若生产已产生用户档案，downgrade 前必须导出相关表，避免不可逆数据丢失。

## 结果与剩余风险

待实现后补充。真实 Telegram/企业微信样本评估和外部端到端测试必须单独标记。
