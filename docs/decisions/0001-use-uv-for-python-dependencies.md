# ADR-0001：使用 uv 管理 Python 与依赖

> 状态：accepted · 日期：2026-07-10

## 背景

后端原先只有手写、精确 pin 的 `requirements.txt`。Python 版本、直接依赖、开发工具和解析后的
传递依赖没有一个统一的项目模型；本地、CI 与 Docker 分别创建环境，容易出现 Python 版本不兼容、
依赖解析漂移和“本地可用但镜像不同”的情况。

项目主要由编码代理维护，需要一个可由仓库机械验证、跨本地/CI/Docker 复现且能明确区分运行时
与开发依赖的入口。

## 决策

- 使用 uv project 模式作为后端唯一 Python 工具链入口。
- `backend/pyproject.toml` 声明 Python 范围、运行时依赖和 `dependency-groups.dev`；迁移时保留
  原直接依赖的精确版本，不把工具迁移伪装成依赖升级。
- `backend/uv.lock` 记录跨平台完整解析并提交版本控制。
- 本地使用 `uv sync --locked --dev` / `uv run --locked`；CI 固定 uv 与 Python 版本；Docker
  固定官方 uv 镜像版本并使用 `uv sync --locked --no-dev`。
- 删除手写 `requirements.txt`，避免两个依赖真相源。依赖变更通过 `uv add`、`uv remove` 或
  `uv lock --upgrade-package` 完成，并同时提交项目元数据与锁文件。

## 备选方案

- **继续 pip + requirements.txt**：改动最少，但不能同时表达 Python、开发依赖与跨平台传递锁；
  需要额外的 compile/export 工具才能获得同等级复现性。
- **Poetry/PDM**：也能提供项目与锁管理，但会引入另一套命令、锁格式和镜像安装路径；当前没有
  足以抵消迁移成本的发布/打包需求。
- **只在本地使用 uv pip**：安装更快，但仍以 requirements 为中心，无法建立统一 project/lock
  工作流。

## 后果

正面影响：开发、CI、Docker 使用同一解析；uv 可自动选择受支持 Python；运行与开发依赖分组；
锁文件漂移能在 `--locked` 下立即失败。

成本与约束：贡献者和 CI 需要受支持的 uv；`uv.lock` 体积增加；更新直接依赖时必须审查锁文件；
官方 uv 版本升级属于工具链变更，需要同时验证本地、CI 和 Docker。

## 验证与复审

- `make backend-check` 必须从 locked sync 开始。
- `scripts/harness_check.py` 要求 pyproject/lock 存在、禁止第二份 requirements，并检查 Makefile、
  CI 与 Docker 都使用 locked sync。
- `uv lock --project backend --check` 和实际后端 Docker build 是发布前验证证据。
- 当项目需要发布 Python package、引入 uv workspace，或 uv 版本范围无法覆盖受支持环境时复审。
