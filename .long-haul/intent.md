# Intent — olap-sql

## 做什么

**用户触发全面刷新。** 项目自 2026-05-14 起处于维护模式（Phase 3 完成），用户 2026-06-22 主动触发新一轮推进。三项核心工作：

### 前提：手动关闭 Issue #14
用户已承诺手动关闭 [Issue #14](https://github.com/AWaterColorPen/olap-sql/issues/14)（yaml.v2 CVE），关闭后本 intent 正式启动。留言模板：
> PR #25 已将 testify 升级至 v1.11.1，gopkg.in/yaml.v2 已从 go.sum 完全消失，CVE-2019-11254 已修复。关闭此 Issue。

### 1. 依赖全面升级
- **Go directive**: 1.24.1 → 1.27（目标最新稳定版，1.27 即将发布时跟进）
- **直接依赖扫描 + 升级**：
  - gorm v1.31.1 → 最新
  - gorm/driver/clickhouse v0.7.0 → 最新
  - gorm/driver/mysql v1.6.0 → 最新
  - gorm/driver/postgres v1.6.0 → 最新
  - gorm/driver/sqlite v1.6.0 → 最新
  - testify v1.11.1 → 最新
  - BurntSushi/toml v1.6.0 → 最新
- **注意**: go-sql-driver/mysql v1.8.1 标志性落后；clickhouse-go/v2 跨度大（上次升级 v2.3.0→v2.45.0），需评估最新版本
- **原则**: 非直接依赖由 `go mod tidy` 自动处理；升级后全量 `go test ./...` 验证

### 2. SDD 建设（从零开始）
当前 `docs/sdd/` 目录不存在，SDD 完全空白。

需建设以下 SDD 文档：
- `schema-mapping.md` — 字典/列映射与 Schema 关系（核心设计）
- `splitter-engine.md` — 分词/拆解引擎（dictionary_splitter.go，最大文件 11KB）
- `translator-chain.md` — 翻译链设计（dictionary_translator.go）
- `dependency-resolution.md` — 依赖图解析（dependency_graph.go + manager.go）
- `adapter-pattern.md` — 适配器模式（dictionary_adapter.go）
- `docs/sdd/README.md` — SDD 索引

### 3. 文档重建与升级
已有文档（docs/，Phase 2 产物）：
- `getting-started.md` / `api.md` / `architecture.md` / `configuration.md` / `examples.md` / `query.md` / `result.md`

需做的事：
- 所有文档与当前 Phase 3 代码库一致性审查（距上次更新约 2 个月）
- README.md 更新：反映当前版本、Go 版本、覆盖率 badge
- CI badge 验证是否仍有效
- 考虑从 `api.md` 拆分出更细粒度的模块文档

## 不做什么

- **不做新功能** — 仅维护性升级和文档建设
- **不做 breaking API 变更** — 接口签名保持不变
- **不做 SDD 之外的超规格文档** — YAGNI
- **不依赖 Issue #14 关闭后立即启动** — 如果用户尚未关闭，先完成其他项并提醒

## 为什么做

olap-sql 在 Phase 3 完成后进入维护模式已 39 天无代码变动。用户决定在维护节奏内做一次全面刷新：
- 依赖随时间积累版本漂移，mysql driver 已落后数个版本
- SDD 长期空白，团队（含 agent）对核心设计缺乏结构化理解
- 文档需要与当前代码同步，避免"文档腐化"

## 成功判据

- [ ] Issue #14 已手动关闭
- [ ] Go directive 升级至合理的当前版本（1.24 最新 patch 或 1.26.x）
- [ ] 所有直接依赖扫描并升级至当前最新稳定版
- [ ] `go test ./...` 全部通过，覆盖率不退化
- [ ] SDD 文档 5 篇 + README.md 索引，覆盖核心模块
- [ ] 所有 docs/ 文档与代码一致性审查完成，过时引用已更新
- [ ] README.md version/badge 信息准确

## 关联

- lifecycle: active
- cadence: 推进期间提升为 daily（完成 3 项工作后恢复 weekly）
- 前序: Phase 3 完成（维护模式，2026-05-14）
- 阻塞: Issue #14 手动关闭（用户侧操作）
