# intent.md — olap-sql 长期迭代立意

---

## 【2026-05-02】当前立意（Phase 3 第三步 — api/models 和 api/types 子包测试覆盖率提升）

**核心目标：** 为 `api/models` 和 `api/types` 两个子包补充基础测试，使其覆盖率从 0% 提升到有意义的水平（目标 60%+）。

### 待执行

1. 新建分支 `feature/phase3-subpkg-coverage`
2. 为 `api/types` 补测（Filter.Expression、Column.GetExpression、FilterOperatorType.IsTree 等核心逻辑）
3. 为 `api/models` 补测（Graph.GetTree、DataSource.IsValid/IsFact/IsDimension、JoinPair/DimensionJoins 等）
4. 全量测试通过后提 PR

### 不做

- 不修改生产代码逻辑
- 不等待 PR #25 合入（此工作独立进行）
- 不追求 100% 覆盖率，只覆盖核心可测逻辑

### 成功判据

`api/types` 和 `api/models` 子包覆盖率均从 0% 提升到 60%+，全量测试通过，PR 已提交。

---

## 【2026-04-30】历史立意（Phase 3 第二步 — 等待 PR #25 合入）

**核心目标：** 等待 PR #25 合入 main。合入后 Phase 3 依赖升级完成，可评估是否进入 Phase 3 第三步或收尾。

### 待执行

1. 确认 PR #25 已合入 main
2. 合入后确认 Phase 3 完成度（覆盖率 82.8% + 全依赖升级均已完成）
3. 可选：进一步提升 api/models 和 api/types 子包的测试覆盖率（目前均为 0%）

### 不做

- 不在 PR 合入前做其他大改动

### 成功判据

PR #25 合入 main，所有测试通过，Phase 3 核心任务（覆盖率 + 依赖升级）均完成。

---

## 【2026-04-28】历史立意（Phase 3 第二步准备 — 依赖升级 — 已推进）✅

**核心目标：** PR #24 合入后，启动 Phase 3 第二步：主要依赖全面升级。

### 待执行

1. 确认 PR #24 已合入 main
2. 新建分支 `feature/phase3-deps-upgrade-round1`，分批升级依赖
   - **Round 1（低风险）**：BurntSushi/toml、cenkalti/backoff、go-faster/* 等工具依赖
   - **Round 2（中风险）**：gorm.io/gorm v1.23.10 → v1.31.1 + 全部 gorm drivers
   - **Round 3（高风险/独立 PR）**：clickhouse-go/v2 v2.3.0 → v2.45.0（先查 CHANGELOG 评估影响）
3. 每次升级后全量运行 `go test ./...`，通过后才提 PR

### 不做

- 不在本阶段修改生产代码逻辑
- 不升级 Go 版本（已是 1.24，够用）
- 不一次升所有依赖（避免大爆炸回归）

### 成功判据

gorm 核心依赖和主要 driver 均升级到当前稳定版，全量测试通过，分 PR 合入 main。

---

## 【2026-04-26】历史立意（Phase 3 第一步 — 已完成）✅

**核心目标：** 将主包测试覆盖率从 79.1% 提升到 80%+（已完成首步：82.8%）。

### 已完成

1. **新增 `coverage_boost_test.go`** — 为 `Clients.SetLogger/BuildSQL`、`Manager.SetLogger/BuildSQL`、`NewTranslator` (direct-SQL 路径)、`FileAdapter.GetMetricsBySource/GetDimensionsBySource`、`DBOption.NewDB` 不支持类型错误路径 添加测试 ✅
2. **覆盖率：79.1% → 82.8%**（超过 80% 目标）✅
3. **PR #24 代码审查意见已处理**（commit 56d672f）— 移除 `TestClients_BuildSQL` 死代码 MockLoad，加强 `TestNewTranslator_DirectSQL` 类型断言 ✅
4. **PR #24 等待合入** ⏳

---

## 【2026-04-24】历史立意（Phase 2 收尾 — 已完成）✅

**核心目标：** 完成 Phase 2 文档建设的收尾任务，为进入 Phase 3 做准备。✅ **已完成**

### 已完成

1. **完善 CONTRIBUTING.md** — 从 stub 扩展为完整贡献指南 ✅
2. **新增 docs/api.md** — 覆盖 Manager、Configuration、Query、Result 的完整 API 参考 ✅
3. **PR #23 已合入 main**（squash merge，commit fa3bd9d）✅

### 下一步

**Phase 2 已全部完成。** 待用户确认后进入 Phase 3：依赖全面更新、测试覆盖率 80%+。

---

## 项目背景

olap-sql 是 Go OLAP 查询 SQL 生成库，现代化迭代计划分三阶段：
- Phase 1（地基）：✅ 已完成 — Go 升级、依赖更新、代码现代化、godoc
- Phase 2（文档建设）：✅ 已完成 — Getting Started/Examples/Architecture/API Reference/CONTRIBUTING 均已完成
- Phase 3（现代化深化）：⬜ 待启动 — 依赖全面更新、测试覆盖率 80%+
