# 处理 harness review 阻断项

> 状态：completed · Owner：Codex · 创建/完成：2026-07-10

## 目标与用户价值

让 harness 分支在不回退生产 OAuth、迁移、时区与部署修复的前提下可安全合并，并把这次暴露的
关键不变量升级为机械检查。

## 非目标

- 不重写现有 OAuth 或部署实现。
- 不拆分已经发布的分支历史为多个 PR；uv 决策通过独立 ADR 稳定记录。
- 不在没有明确推送指令时改写远端分支。

## 背景与原有行为

review 发现分支基于 `67e5da0`，落后 `origin/main` 的 4 个生产修复，同时 harness 只检查内部
`app.*` 依赖，未覆盖 Alembic 单链、SQLModel 时区列和 domain 外部框架导入。

## 验收结果

- [x] 分支 rebase 到 `ed869bd`，保留 `202607100002_repair_auth_schema.py`。
- [x] 保留 timezone-aware SQLModel 字段和分阶段、强制重跑 migrate 的部署流程。
- [x] harness 拒绝缺失/重复/多 head 的 Alembic revision。
- [x] harness 拒绝未显式配置 timezone 的持久化 datetime 字段。
- [x] harness 拒绝 domain 导入框架、数据库、队列和 provider 实现包。
- [x] uv 决策有 accepted ADR，并从 ADR 索引可达。
- [x] 后端、前端、workflow、Compose、Dockerfile 和 release 分支契约验证通过。

## 影响面与风险

修改 harness 脚本、知识库和分支历史，不改变生产业务代码。rebase 后本地与旧远端历史分叉，后续
发布必须先核对远端 SHA，并使用 `git push --force-with-lease`，不能使用无保护的 force push。

## 实施记录

- [x] fetch/rebase 并核对生产修复。
- [x] 实现三类机械检查与 5 个回归测试。
- [x] 新增 uv ADR。
- [x] 完整检查并归档计划。

## 进度与发现日志

- 2026-07-10：`git rebase origin/main` 无冲突完成；repair migration、timezone 字段与分阶段
  migrate/runtime 启动均保留。
- 2026-07-10：现有 `workflow_run` release 触发改造在 rebase 后也保留。
- 2026-07-10：`make check` 覆盖 harness、后端和前端并完整通过。

## 决策日志

- 2026-07-10：不拆分已发布提交；通过 rebase + accepted ADR 解决基线与工具链决策问题。
- 2026-07-10：迁移检查以 revision graph 单根单 head + repair migration 必需文件为当前保护层；
  数据库实际 upgrade 仍由 Alembic/集成环境验证。
- 2026-07-10：domain 暂时允许 `pydantic`，因为现有端口契约使用它；阻止 Web、ORM、队列、
  provider 和基础设施实现包。

## 验证记录

| 命令/场景 | 结果 |
| --- | --- |
| `make harness-check` | 通过，文档、边界、不变量和回归测试均通过 |
| `make check` | 通过；后端 12 tests，前端 lint/typecheck/build 通过 |
| `uv run --locked alembic heads` | `202607100002 (head)` |
| SQLAlchemy metadata timezone 检查 | 22 个 datetime column 均为 timezone-aware |
| workflow YAML + release 分支契约 | 3 个 workflow 均通过 |
| `docker compose config --quiet` | 通过 |
| `docker build --check -f Dockerfile .` | 通过，无警告 |
| `git diff --check` | 通过 |

## 回滚与恢复

远端仍保留 rebase 前的旧提交。若检查实现有误，可只回退本次 review-fix 变更，不回退
`origin/main` 的生产修复。

## 结果与剩余风险

5 条 review finding 已处理。剩余操作仅为提交并用 lease 更新远端；这是有意保留的发布边界，
需在明确要求提交/推送后执行。
