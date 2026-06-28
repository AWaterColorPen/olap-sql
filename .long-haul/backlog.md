# Backlog — olap-sql

## P0（阻塞当前方向）

- [ ] gh#14 | 关闭 yaml.v2 CVE 相关 Issue | external:github#14 | PR #25 已移除 yaml.v2，需 owner 手动关闭 Issue #14 并说明修复原因

## P1（高价值近期做）

- [ ] BK-2026-06-19-1 | 配置 Go 运行以执行依赖扫描 | capability-gap | 本地无 Go 运行时，无法执行 `go list -m -u all`；安装后可在 weekly 维护中扫描过期依赖

## P2（可做可不做）

- [ ] BK-2026-06-19-2 | api/types/clause.go 集成测试 | reflect:2026-06-13 | 需要真实 DB 连接，当前暂不做
- [ ] BK-2026-06-19-3 | 性能基准测试 | reflect:2026-06-13 | 维护期候选方向
- [ ] BK-2026-06-19-4 | CI/CD 改进（覆盖率上报、Lint Action） | reflect:2026-06-13 | 维护期候选方向

## 已关闭

- [x] gh#25 | 依赖全面升级 | done 2026-05-14 | 已合入 main，移除 yaml.v2
- [x] gh#26 | api/models + api/types 子包覆盖率提升 | done 2026-05-14 | 已合入 main
