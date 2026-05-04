# reflect.md — olap-sql 回看记录

---

## 【2026-05-04】PR 等待期间的阶段性回看

### 我们在哪？

PR #25（依赖全面升级）和 PR #26（子包覆盖率提升）仍处于 open 状态，均可直接合入（mergeable=true），CI 全部通过（build: success），无审查意见，纯等待 owner 操作。

当前代码库状态（main 分支）：
- 主包覆盖率：82.8%
- api/models 覆盖率：0%（PR #26 合入后将为 94.6%）
- api/types 覆盖率：0%（PR #26 合入后将为 60.3%）
- 依赖：尚未升级（PR #25 合入后全面升级到 2026 年稳定版）

**新发现**：Issue #14 报告了 gopkg.in/yaml.v2@v2.2.2 的 CVE-2019-11254 安全漏洞（拒绝服务）。经排查，该依赖是 testify@v1.4.0/v1.5.1 的传递依赖，项目本身不直接使用 yaml.v2。**PR #25 将 testify 升级到 v1.11.1，go.sum 中 yaml.v2 条目完全消失，等价于修复了此漏洞。** 因此 PR #25 不仅是依赖现代化，也是一个安全修复。

### 我们离目标更近了吗？

是的，Phase 3 三步核心工作已全部完成并有 PR，等待合入：
- PR #24 ✅ 已合入（覆盖率 82.8%）
- PR #25 ⏳ open，可合入，CI 通过（修复 CVE-2019-11254 + 依赖现代化）
- PR #26 ⏳ open，可合入，CI 通过（子包覆盖率提升）

Phase 3 的完成只差 owner 点击 merge。

### 下一步最应该做什么？

1. **无主动操作可做** — 两个 PR 均无 review comments，只能等 owner。不需要发 PR comment 催合入。
2. **PR #25 合入后，关闭 Issue #14** — PR #25 间接修复了 yaml.v2 CVE，合入后可在 Issue #14 中留一条说明并关闭。
3. **两个 PR 均合入后，进行 Phase 3 完整收尾回看**，评估：
   - 项目是否可以宣告进入「维护模式」
   - api/types/clause.go（BuildSQL/BuildDB）的 0% 覆盖是否需要单独处理（需要集成测试，当前不强求）
   - 是否有新的方向性需求

### 风险/注意事项

- PR #25 和 PR #26 已提交超过 2 天无动作，这是 owner 操作节奏问题，不是项目问题
- PR #25 合入时如果 main 分支有新提交，可能需要 rebase；目前 main 无新 commit，直接合入无冲突
- yaml.v2 漏洞（CVE-2019-11254）在当前 main 仍存在，但 PR #25 合入后即消失，优先级可以接受

### 建议

下次唤醒时：
1. 先检查 PR #25 和 PR #26 是否合入
2. 若均已合入：进行 Phase 3 完整收尾回看，在 Issue #14 中说明 yaml.v2 漏洞已通过依赖升级修复
3. 若仍未合入：继续等待，无新操作需求（除非出现审查意见）
4. 若 PR 等待超过 1 周无回应：可考虑通过其他渠道提醒 owner，或在 PR 上补充说明安全修复的重要性

---

## 【2026-05-02】Phase 3 子包覆盖率提升后回看

### 我们在哪？

Phase 3 三步核心工作均已完成并提 PR：
- PR #24（覆盖率提升到 82.8%）— 已合入 ✅
- PR #25（依赖全面升级）— open，等待仓库 owner ⏳
- PR #26（api/models+api/types 子包覆盖率）— 刚提交 ⏳

当前覆盖率：main=82.8%，api/models=94.6%，api/types=60.3%（均从 0% 提升）

### 我们离目标更近了吗？

是的。Phase 3 原定两个目标（覆盖率 80%+、依赖全面升级）均已达成。本次额外推进了子包覆盖率，整体测试质量进一步提升。

### 下一步最应该做什么？

1. **等待 PR #25 和 PR #26 合入** — 无需主动操作，等仓库 owner
2. **两个 PR 均合入后**，进行一次 Phase 3 完整收尾回看，评估：
   - Phase 3 是否全部完成
   - api/types 剩余 39.7% 未覆盖部分（clause.go 需要 gorm DB，属于集成测试范畴，当前暂时不做）
   - 是否有必要进入新的迭代阶段
3. **Phase 3 收尾后**可考虑宣告项目进入维护模式，或开始讨论新方向

### 风险/注意事项

- PR #25 中 clickhouse-go/v2 v2.45.0 vs gorm/driver/clickhouse v0.7.0 的潜在不兼容问题仍存在，但测试通过
- api/types/clause.go 的 BuildSQL/BuildDB 覆盖率为 0%，这些函数是核心 SQL 构建逻辑，但需要真实 DB 连接（gorm 集成）才能测试

### 建议

下次唤醒时：
1. 先检查 PR #25 和 PR #26 是否合入
2. 若均已合入，进行 Phase 3 完整收尾回看
3. 若仍有未合入的 PR，检查是否有审查意见需处理

---

## 【2026-04-30】Phase 3 依赖全面升级完成后回看

### 我们在哪？

Phase 3 第二步（依赖全面升级）已完成。PR #25 提交 `feature/phase3-deps-upgrade-round1`，包含所有主要依赖升级：gorm v1.31.1、所有 gorm drivers v1.6.0、clickhouse-go/v2 v2.45.0、otel v1.41.0 等。全量测试通过，覆盖率保持 82.8%。PR 等待仓库 owner 合入。

### 我们离目标更近了吗？

是的。Phase 3 两个核心目标（覆盖率 80%+、依赖全面升级）均已完成并提 PR。两个 PR 都已推送，等待合入：
- PR #24（覆盖率提升到 82.8%）— 已合入 ✅
- PR #25（依赖全面升级）— 等待合入 ⏳

### 下一步最应该做什么？

1. **等待 PR #25 合入**
2. **合入后评估 Phase 3 是否完整收尾**，或考虑：
   - 对 api/models 和 api/types 子包补基础测试（提升 total 覆盖率）
   - 评估是否有其他 Phase 3 目标未完成
3. Phase 3 完成后，可以做一次完整的项目回看，评估是否进入新的迭代阶段

### 风险/注意事项

- clickhouse-go/v2 v2.45.0 比 gorm/driver/clickhouse v0.7.0 要新（driver 依赖 v2.30.0），直接升到 v2.45.0 后测试通过，但 driver 层未完全验证所有 ClickHouse 功能路径
- golang.org/x/sync 经 `go mod tidy` 的传递升级，go.mod 中 `go 1.24` 被 clickhouse-go 改成了 `go 1.24.1 + toolchain go1.24.11`，这是标准行为，无需担心

### 建议

下次唤醒时：
1. 先检查 PR #25 是否合入
2. 若已合入，做一次 Phase 3 完整回看，评估项目整体健康度
3. 可考虑为 api/models 和 api/types 补充基础测试，提升 total 覆盖率

---

## 【2026-04-28】PR #24 审查处理后回看

### 我们在哪？

Phase 3 第一步（覆盖率提升）的 PR #24 已处理完全部审查意见（commit 56d672f 已推送），等待仓库 owner 合入。依赖升级信息已初步扫描：gorm.io/gorm v1.23.10 → v1.31.1，gorm 各 driver v1.3.6 → v1.6.0，clickhouse-go/v2 v2.3.0 → v2.45.0，mysql driver v1.6.0 → v1.9.3。

### 我们离目标更近了吗？

是的。PR #24 的代码质量问题已修复，覆盖率稳定在 82.8%。Phase 3 第一步只差合入。

### 下一步最应该做什么？

1. **确认 PR #24 合入状态** — 下次唤醒优先检查
2. **若已合入，启动 Phase 3 第二步：依赖升级**
   - 分优先级：低风险先升（BurntSushi/toml、testify 等工具依赖）→ 中风险（gorm 核心 + drivers）→ 高风险最后（clickhouse-go/v2，跨 minor 版本跨度大）
   - 每次升一批，独立 PR，避免大爆炸

### 风险/注意事项

- clickhouse-go/v2 v2.3.0 → v2.45.0 跨度极大（42个 minor），需先查 CHANGELOG，评估 breaking change
- gorm v1.23.10 → v1.31.1 也有较大跨度，关注 API 变化
- 升级后必须全量运行 `go test ./...` 确认无回归

### 建议

下次唤醒时：
1. 先检查 PR #24 是否合入
2. 若已合入，拉取最新 main，新建分支 feature/phase3-deps-upgrade-round1，开始低风险依赖升级
3. 若未合入，查看是否有新的审查意见需处理

---

## 【2026-04-26】Phase 3 第一步完成后回看

### 我们在哪？

Phase 3 第一步（测试覆盖率提升）已完成。主包覆盖率已从 79.1% 提升到 82.8%，超过 80% 目标。PR #24 已提交，等待合入。

### 我们离目标更近了吗？

是的。Phase 3 第一步的成功判据已达成。但整体 total（含 api/models 和 api/types 子包）约 52.9%，这两个子包完全没有测试，如果未来要进一步提升整体覆盖率，需要对子包加测。

### 下一步最应该做什么？

1. 合入 PR #24（feature/phase3-coverage-boost）
2. 评估主要依赖升级（gorm、clickhouse-go）——这是 Phase 3 第二步
3. 若要进一步提升 total 覆盖率，可考虑对 api/models 和 api/types 子包补几个基础测试

### 风险/注意事项

- 依赖升级可能引入 breaking change，需要单独 PR + CHANGELOG
- clickhouse-go 从 v2 → v3 可能有较大 API 变动，需要先评估影响再动手

### 建议

下次唤醒时：
1. 先确认 PR #24 是否已合入
2. 若已合入，开始评估依赖升级：`go list -m -u all` 查看可用更新
3. 优先升级风险较低的依赖（非 DB 驱动层）

---

## 【2026-04-24】Phase 2 完成后回看

### 我们在哪？

Phase 2（文档建设）已全部完成。项目现在有完整的使用者文档（Getting Started、Configuration、Query、Result、Examples、API Reference）和贡献者文档（Architecture、CONTRIBUTING）。代码库本身处于良好状态：Go 1.24、现代化语法、通过全部测试。

### 我们离目标更近了吗？

是的。按 spec 定义的三阶段计划：
- Phase 1 ✅
- Phase 2 ✅（本次完成了最后两个缺口：API Reference + CONTRIBUTING）
- Phase 3 ⬜ 待启动

### 下一步最应该做什么？

等待用户确认进入 Phase 3。Phase 3 的主要工作：
1. **主要依赖全面升级** — 评估 gorm、clickhouse-go 等的最新版本，分 PR 升级
2. **提升测试覆盖率到 80%+** — 重点补充 `dictionary_translator.go` 和 `dependency_graph.go` 的单元测试

### 风险/注意事项

- Phase 3 的依赖升级可能引入 breaking change，需要单独 PR + CHANGELOG
- clickhouse-go 从 v2 → v3 可能有较大 API 变动，需要先评估影响再动手
- 测试覆盖率提升前，应先了解现有覆盖率基线（`go test -cover ./...`）

### 建议

下次唤醒时：
1. 先向用户确认是否可以进入 Phase 3
2. 若确认，先运行 `go test -cover ./...` 获取覆盖率基线
3. 再评估依赖最新版本，制定升级优先级
