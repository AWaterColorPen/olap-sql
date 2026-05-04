# readiness.md — olap-sql 环境就绪性清单

## 语言与运行时

- **Go 1.24+** required (project uses range-over-integer, min/max builtins)
- 检查：`go version`

## 依赖安装

```bash
go mod tidy
```

## 测试验证

```bash
go test ./...
```

（依赖本地 DB 的集成测试可能需要 ClickHouse / MySQL / PostgreSQL / SQLite，但 SQLite 测试应开箱可用）

## 分支策略

- **主分支：** `main`（不直接 push，通过 feature 分支 + PR 合入）
- 每次改动单独一个 feature 分支 → PR → merge
- Breaking change 需单独 PR + CHANGELOG 条目

## Git remote 格式

```
https://AWaterColorPen:<token>@github.com/AWaterColorPen/olap-sql.git
```

Token 由 cron 任务指令提供，不持久化到文件。

## 每次唤醒复核项

1. `git status` 干净（无未提交改动）
2. 当前在 `main` 分支（除非有进行中的 feature 分支）
3. `go test ./...` 通过（sqlite 相关测试）
4. 检查 `.long-haul/intent.md` 中的当前立意，判断继续方向
