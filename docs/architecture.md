# Architecture

This document describes the internal design of olap-sql for contributors who want to understand, extend, or debug the library.

---

## Table of Contents

- [High-level overview](#high-level-overview)
- [Component responsibilities](#component-responsibilities)
  - [Configuration & Manager](#configuration--manager)
  - [Dictionary & Adapter](#dictionary--adapter)
  - [Translator](#translator)
  - [Clause & Statement](#clause--statement)
  - [Clients & Database layer](#clients--database-layer)
  - [Result](#result)
- [Request lifecycle](#request-lifecycle)
- [Data-flow diagram](#data-flow-diagram)
- [Package layout](#package-layout)
- [Key design decisions](#key-design-decisions)
- [Extension points](#extension-points)

---

## High-level overview

olap-sql sits between your application code and your OLAP database.  
You describe your data model once in a TOML schema file. After that, your code only works with high-level query objects — no SQL strings, no JOIN logic, no database-specific dialect.

```
Your code
  │  types.Query
  ▼
Manager ──► Dictionary ──► Translator ──► types.Clause
                                              │  SQL string / *gorm.DB
                                              ▼
                                          Clients ──► Database
                                              │  []map[string]any rows
                                              ▼
                                          types.Result
```

---

## Component responsibilities

### Configuration & Manager

`Configuration` is the single struct you pass to `NewManager`.  It carries two independent sections:

| Field              | Purpose |
|--------------------|---------|
| `ClientsOption`    | Maps string keys → `DBOption` (DSN + type). Drives how GORM connections are opened. |
| `DictionaryOption` | Points to the TOML schema file (or another adapter). Drives schema loading. |

`Manager` is the public entry point.  It owns a `Clients` map and a `Dictionary`, and exposes four methods:

| Method               | Description |
|----------------------|-------------|
| `BuildSQL(query)`    | Translate query → SQL string (dry-run, no DB call) |
| `BuildTransaction(q)`| Translate query → `*gorm.DB` ready to execute |
| `RunSync(query)`     | Execute and return all rows in a slice |
| `RunChan(query)`     | Execute and stream rows over a channel |

### Dictionary & Adapter

`Dictionary` wraps an `IAdapter`.  The adapter is responsible for loading and serving the virtual schema: sets, sources, metrics, and dimensions.

```
Dictionary
  └── IAdapter (interface)
        └── fileAdapter  ← current implementation; loads from a TOML file
```

`IAdapter` exposes:

- `GetDataSetByKey(name)` — returns the named `Set` with its configured DB type and data source
- `GetSourceByKey(name)` — returns a `Source` (table or join definition)
- `GetMetricsByKey(name)` — returns a `Metric` definition
- `GetDimensionByKey(name)` — returns a `Dimension` definition

New adapter types (e.g. load schema from a database, HTTP endpoint, or in-memory struct) can be added by implementing `IAdapter`.

### Translator

`Translator` converts a `types.Query` into a `types.Clause`.  This is where the core query-building logic lives.

The translation steps are:

1. **Set resolution** — find the `Set` named by `query.DataSetName`; determine the target DB type
2. **Metric expansion** — for each metric name, resolve its definition (and recursively resolve composition dependencies)
3. **Dimension expansion** — same for dimensions
4. **Filter translation** — convert `Filter` structs to SQL `WHERE` fragments; handle tree operators (`AND`, `OR`) recursively
5. **ORDER BY / LIMIT** — append ordering and pagination clauses
6. **JOIN resolution** — if the source is a `fact_dimension_join` or `merge_join`, synthesise the necessary `JOIN ON` expressions
7. **Clause assembly** — package everything into a `Clause` that knows its DB type, dataset, and all SQL fragments

The translator is created via `NewTranslator(*TranslatorOption)` and is an internal type; callers use `Dictionary.Translate(query)`.

### Clause & Statement

`types.Clause` is the output of the Translator.  It is a database-backend-specific object that wraps the translated SQL fragments and knows how to apply them to a `*gorm.DB`.

```go
type Clause interface {
    GetDBType() DBType
    GetDataset() string
    BuildDB(db *gorm.DB) (*gorm.DB, error)
    BuildSQL(db *gorm.DB) (string, error)
}
```

Each SQL fragment (SELECT column, WHERE condition, ORDER BY, etc.) implements the `Statement` interface:

```go
type Statement interface {
    Expression() (string, error)
    Alias() (string, error)
    Statement() (string, error)
}
```

This makes it straightforward to add new metric types or filter operators: implement `Statement`, register the type constant, and wire it in the translator.

### Clients & Database layer

`Clients` is a `map[string]*gorm.DB` keyed by `"<dbtype>"` or `"<dbtype>/<dataset>"`.

- `RegisterByOption` opens connections from `ClientsOption` at startup
- `Get(dbType, dataset)` performs a two-level lookup (dataset-specific → type-level fallback)
- `BuildDB(clause)` and `BuildSQL(clause)` are convenience wrappers that pick the right connection and execute the clause

GORM is used as the SQL builder and connection pool.  Supported drivers: ClickHouse, MySQL, PostgreSQL, SQLite.

### Result

`types.Result` is the output handed back to callers:

```go
type Result struct {
    Dimensions []string         // ordered list of column names (dimensions + metrics)
    Source     []map[string]any // one map per row; keys are column names
}
```

`SetDimensions(query)` copies `query.Dimensions` then appends `query.Metrics`, preserving the caller's column order.  Each row in `Source` maps a column name to its raw value as returned by the GORM scanner.

---

## Request lifecycle

```
manager.RunSync(query)
  │
  ├─ query.TranslateTimeIntervalToFilter()   // expand TimeInterval → Filter
  │
  ├─ dictionary.Translate(query)
  │     ├─ adapter.GetDataSetByKey(...)
  │     ├─ NewTranslator(option)
  │     └─ translator.Translate(query)
  │           ├─ resolve metrics + dimensions
  │           ├─ translate filters recursively
  │           ├─ build JOIN chain (if needed)
  │           └─ assemble Clause
  │
  ├─ clients.BuildDB(clause)
  │     ├─ clients.Get(dbType, dataset)      // pick the right *gorm.DB
  │     └─ clause.BuildDB(db)               // apply SQL fragments to gorm.DB
  │
  ├─ RunSync(db)                            // db.Scan(&rows)
  │
  └─ BuildResultSync(query, rows)           // wrap rows in Result
```

---

## Data-flow diagram

```
 ┌─────────────────────────────────────────────────────────────┐
 │                        Your application                      │
 │                                                              │
 │  query := &types.Query{DataSetName:"wikistat", ...}          │
 │  result, err := manager.RunSync(query)                       │
 └────────────────────────────┬─────────────────────────────────┘
                               │  *types.Query
                               ▼
 ┌─────────────────────────────────────────────────────────────┐
 │                         Manager                              │
 │                                                              │
 │  1. TranslateTimeIntervalToFilter()                          │
 │  2. dictionary.Translate(query) ──────────────────────────► │
 │                                    Dictionary                 │
 │                                      └─ IAdapter             │
 │                                           └─ fileAdapter     │
 │                                                (TOML)        │
 │  3. clients.BuildDB(clause)                                  │
 │  4. RunSync / RunChan                                        │
 │  5. BuildResult*                                             │
 └────────────────────────────┬─────────────────────────────────┘
                               │  *types.Result
                               ▼
 ┌─────────────────────────────────────────────────────────────┐
 │  result.Dimensions  []string          (column names)         │
 │  result.Source      []map[string]any  (rows)                 │
 └─────────────────────────────────────────────────────────────┘
```

---

## Package layout

```
olap-sql/
├── *.go               # Public API: Manager, Configuration, Clients, Dictionary, ...
├── api/
│   └── types/         # Shared type definitions (Query, Filter, Result, Clause, ...)
├── test/
│   ├── dictionary.ck.toml      # ClickHouse test schema
│   └── dictionary.sqlite.toml  # SQLite test schema
├── docs/              # Documentation (you are here)
└── scripts/           # Utility shell scripts
```

| File / package | Responsibility |
|----------------|----------------|
| `manager.go`   | Public `Manager` type; `RunSync`, `RunChan`, `BuildSQL`, `BuildTransaction` |
| `client.go`    | `Clients` map; `RegisterByOption`, `Get`, `BuildDB`, `BuildSQL` |
| `configuration.go` | `Configuration` struct |
| `database.go`  | `DBOption`; opens GORM connections by `DBType` |
| `dictionary.go`| `Dictionary`; wraps `IAdapter`; exposes `Translate` |
| `dictionary_adapter.go` | `IAdapter` interface; `fileAdapter` implementation |
| `dictionary_column.go`  | Column-level schema objects (metric/dimension field descriptors) |
| `dictionary_translator.go` | `Translator`; converts `Query` → `Clause` |
| `dictionary_splitter.go`   | Splits a joined source into its constituent tables and JOIN keys |
| `run.go`       | `RunSync`, `RunChan`, `BuildResult*` helpers |
| `dependency_graph.go` | Resolves composition metric/dimension dependency chains |
| `api/types/`   | All shared value types: `Query`, `Filter`, `Clause`, `Result`, `Statement`, ... |

---

## Key design decisions

**Schema-driven, not code-driven.**  
The TOML schema is the single source of truth for which tables exist, how they join, and what each metric means.  Application code only references names (strings).  This means you can change the underlying table structure without touching Go code.

**Composition by name reference.**  
Derived metrics (`METRIC_DIVIDE`, `METRIC_ADD`, etc.) reference their inputs by `"<source>.<metric>"` strings.  The translator resolves these recursively, which allows arbitrarily deep derivation chains.

**GORM as a dialect bridge.**  
Rather than hand-rolling SQL for each database, olap-sql builds a `*gorm.DB` with `Select`, `Where`, `Joins`, `Group`, `Order`, `Limit` calls.  GORM translates these to dialect-appropriate SQL, giving olap-sql ClickHouse / MySQL / Postgres / SQLite support for free.

**Flat result model.**  
`Result.Source` is `[]map[string]any` — no generated structs, no reflection-based mapping.  This trades a small runtime cost for zero code generation and easy JSON serialisation.

---

## Extension points

| What to extend | Where to look |
|----------------|---------------|
| New database backend | Add a case to `getDialect` in `database.go`; add the GORM driver dependency in `go.mod` |
| New metric type | Add a constant in `api/types/metric.go`; handle it in the translator and the `Statement` implementation |
| New dimension type | Same pattern as metric types |
| New filter operator | Add a constant in `api/types/filter.go`; handle it in `Filter.Expression()` |
| Load schema from DB/API | Implement `IAdapter` in `dictionary_adapter.go` |
| Custom SQL fragment | Implement `types.Statement` and return it from the translator |
