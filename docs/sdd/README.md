# olap-sql 软件设计文档（SDD）索引

本目录记录 olap-sql 核心模块的设计决策与实现原理，面向需要深入理解代码结构或进行二次开发的维护者。

## 文档清单

| 文档 | 对应源码 | 核心内容 |
|---|---|---|
| `schema-mapping.md` | `dictionary.go`, `dictionary_column.go`, `api/models/`, `api/types/` | 字典（Dictionary）如何映射业务语义（DataSet / Metric / Dimension）到物理 Schema 与列定义 |
| `splitter-engine.md` | `dictionary_splitter.go` | 查询拆分引擎：按数据源类型（Fact / Join / Merge）将 NormalClause 拆分为子查询 |
| `translator-chain.md` | `dictionary_translator.go` | 翻译链：从 `types.Query` 到 `types.Clause` 的完整转换流程，含 DirectSQL 快捷路径 |
| `dependency-resolution.md` | `dependency_graph.go`, `manager.go` | 依赖图构建与解析：Metric/Dimension 跨表依赖如何被收集、折叠与校验 |
| `adapter-pattern.md` | `dictionary_adapter.go` | 适配器模式：IAdapter 接口、FILE 适配器实现，以及 AdapterOption 配置模型 |

## 阅读顺序

建议按以下顺序阅读，以建立从整体到局部的认知：

1. `adapter-pattern.md` — 理解配置如何被加载为内存中的 Schema 适配器
2. `schema-mapping.md` — 理解 DataSet、DataSource、Metric、Dimension 之间的关系
3. `dependency-resolution.md` — 理解依赖图如何表达跨表引用
4. `translator-chain.md` — 理解查询如何被翻译成 Clause
5. `splitter-engine.md` — 理解复杂查询如何被拆分为可执行的子查询

## 设计原则

- **业务语义与物理实现分离**：上层使用 Metric/Dimension/DataSet 等 OLAP 概念，下层通过 Adapter 映射到具体数据库的表和列
- **声明式查询**：用户通过 `types.Query` 描述想要什么数据，翻译引擎负责生成对应 SQL
- **后端无关的 Clause 中间表示**：`types.Clause` 作为不同数据库方言的统一中间产物
- **依赖图驱动**：所有跨表引用都先被解析为依赖图，再生成 JOIN 或子查询

## 待补充事项

- [ ] `schema-mapping.md` 初稿
- [ ] `splitter-engine.md` 初稿
- [ ] `translator-chain.md` 初稿
- [ ] `dependency-resolution.md` 初稿
- [x] `adapter-pattern.md` 初稿
- [ ] 与 `docs/architecture.md` 的交叉引用审查
