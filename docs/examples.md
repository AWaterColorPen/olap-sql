# Examples

This page shows common real-world usage patterns for olap-sql.  
All examples use **ClickHouse** unless stated otherwise; the query syntax is identical for other backends — only the DSN and a few DB-specific SQL functions differ.

---

## Table of Contents

- [1. Simple single-table query (ClickHouse)](#1-simple-single-table-query-clickhouse)
- [2. Time-range filter with `time_interval`](#2-time-range-filter-with-time_interval)
- [3. Custom metrics formula (division / ratio)](#3-custom-metrics-formula-division--ratio)
- [4. Multi-table JOIN query](#4-multi-table-join-query)
- [5. Synchronous vs streaming (RunSync vs RunChan)](#5-synchronous-vs-streaming-runsync-vs-runchan)
- [6. Debug: inspect generated SQL without running it](#6-debug-inspect-generated-sql-without-running-it)
- [7. SQLite quick-start (no ClickHouse needed)](#7-sqlite-quick-start-no-clickhouse-needed)

---

## 1. Simple single-table query (ClickHouse)

The simplest case: one metric, one dimension, no filters.

**TOML schema** (`olap.toml`):

```toml
sets = [
  {name = "wikistat", type = "clickhouse", data_source = "wikistat"},
]

sources = [
  {database = "", name = "wikistat", type = "fact"},
]

metrics = [
  {data_source = "wikistat", type = "METRIC_SUM", name = "hits", field_name = "hits", value_type = "VALUE_INTEGER"},
]

dimensions = [
  {data_source = "wikistat", type = "DIMENSION_SINGLE", name = "date", field_name = "date", value_type = "VALUE_STRING"},
  {data_source = "wikistat", type = "DIMENSION_SINGLE", name = "project", field_name = "project", value_type = "VALUE_STRING"},
]
```

**Go code**:

```go
package main

import (
    "fmt"
    "log"

    olapsql "github.com/awatercolorpen/olap-sql"
    "github.com/awatercolorpen/olap-sql/api/types"
)

func main() {
    cfg := &olapsql.Configuration{
        ClientsOption: map[string]*olapsql.DBOption{
            "clickhouse": {
                DSN:  "clickhouse://localhost:9000/default",
                Type: types.DBTypeClickHouse,
            },
        },
        DictionaryOption: &olapsql.Option{
            AdapterOption: olapsql.AdapterOption{Dsn: "olap.toml"},
        },
    }

    manager, err := olapsql.NewManager(cfg)
    if err != nil {
        log.Fatal(err)
    }

    query := &types.Query{
        DataSetName: "wikistat",
        Metrics:     []string{"hits"},
        Dimensions:  []string{"project"},
    }

    result, err := manager.RunSync(query)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("columns:", result.Dimensions)
    for _, row := range result.Source {
        fmt.Println(row)
    }
}
```

**Generated SQL**:

```sql
SELECT
    wikistat.project AS project,
    1.0 * SUM(wikistat.hits) AS hits
FROM wikistat AS wikistat
GROUP BY wikistat.project
```

---

## 2. Time-range filter with `time_interval`

`TimeInterval` is a shorthand for a pair of `>=` / `<` filters on a timestamp/date column.  
It is automatically expanded into the `WHERE` clause.

```go
query := &types.Query{
    DataSetName: "wikistat",
    TimeInterval: &types.TimeInterval{
        Name:  "date",           // must match a declared dimension name
        Start: "2021-05-06",     // inclusive
        End:   "2021-05-08",     // exclusive
    },
    Metrics:    []string{"hits"},
    Dimensions: []string{"date"},
}
```

**Generated SQL**:

```sql
SELECT
    wikistat.date AS date,
    1.0 * SUM(wikistat.hits) AS hits
FROM wikistat AS wikistat
WHERE (wikistat.date >= '2021-05-06' AND wikistat.date < '2021-05-08')
GROUP BY wikistat.date
```

You can also add extra filters on top of the time interval:

```go
query := &types.Query{
    DataSetName: "wikistat",
    TimeInterval: &types.TimeInterval{
        Name:  "date",
        Start: "2021-05-01",
        End:   "2021-06-01",
    },
    Metrics:    []string{"hits"},
    Dimensions: []string{"date", "project"},
    Filters: []*types.Filter{
        {
            OperatorType: types.FilterOperatorTypeIn,
            Name:         "project",
            Value:        []any{"en", "de", "fr"},
        },
    },
    Orders: []*types.OrderBy{
        {Name: "hits", Direction: types.OrderDirectionTypeDescending},
    },
    Limit: &types.Limit{Limit: 10},
}
```

**Generated SQL**:

```sql
SELECT
    wikistat.date    AS date,
    wikistat.project AS project,
    1.0 * SUM(wikistat.hits) AS hits
FROM wikistat AS wikistat
WHERE
    (wikistat.date >= '2021-05-01' AND wikistat.date < '2021-06-01')
    AND wikistat.project IN ('en', 'de', 'fr')
GROUP BY wikistat.date, wikistat.project
ORDER BY hits DESC
LIMIT 10
```

---

## 3. Custom metrics formula (division / ratio)

olap-sql supports **composition metrics** — derived metrics built from other metrics.  
A common pattern is a ratio like `hits_avg = hits / count`.

**TOML schema**:

```toml
metrics = [
  # Base metrics
  {data_source = "wikistat", type = "METRIC_SUM",   name = "hits",     field_name = "hits", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_COUNT",  name = "count",    field_name = "*",    value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_SUM",    name = "size_sum", field_name = "size", value_type = "VALUE_INTEGER"},

  # Derived (composition) metrics
  # hits_avg = hits / count  (average hits per row)
  {data_source = "wikistat", type = "METRIC_DIVIDE", name = "hits_avg",     value_type = "VALUE_FLOAT",
   dependency = ["wikistat.hits", "wikistat.count"]},

  # size_per_hit = size_sum / hits
  {data_source = "wikistat", type = "METRIC_DIVIDE", name = "size_per_hit", value_type = "VALUE_FLOAT",
   dependency = ["wikistat.size_sum", "wikistat.hits"]},
]
```

**Query**:

```go
query := &types.Query{
    DataSetName: "wikistat",
    Metrics:     []string{"hits", "hits_avg", "size_per_hit"},
    Dimensions:  []string{"date"},
}
```

**Generated SQL**:

```sql
SELECT
    wikistat.date AS date,
    1.0 * SUM(wikistat.hits) AS hits,
    (1.0 * SUM(wikistat.hits)) / NULLIF(COUNT(*), 0)                   AS hits_avg,
    (1.0 * SUM(wikistat.size)) / NULLIF((1.0 * SUM(wikistat.hits)), 0) AS size_per_hit
FROM wikistat AS wikistat
GROUP BY wikistat.date
```

Other supported composition operators:

| TOML type         | SQL operator   | Example                                     |
|-------------------|----------------|---------------------------------------------|
| `METRIC_ADD`      | `+`            | `dependency = ["wikistat.a", "wikistat.b"]` |
| `METRIC_SUBTRACT` | `-`            | `dependency = ["wikistat.a", "wikistat.b"]` |
| `METRIC_MULTIPLY` | `*`            | `dependency = ["wikistat.a", "wikistat.b"]` |
| `METRIC_DIVIDE`   | `/ NULLIF(,0)` | `dependency = ["wikistat.a", "wikistat.b"]` |

---

## 4. Multi-table JOIN query

olap-sql can join a **fact table** to one or more **dimension tables** automatically.  
You define the join keys in the TOML schema; the generated SQL handles the `JOIN ON`.

### Schema

```toml
sets = [
  {name = "wikistat_join", type = "clickhouse", data_source = "wikistat_base"},
]

sources = [
  # Fact table
  {database = "", name = "wikistat",        type = "fact"},
  # Dimension tables
  {database = "", name = "wikistat_relate",  type = "dimension"},
  {database = "", name = "wikistat_class",   type = "dimension"},
  # Virtual joined source: wikistat ⟶ wikistat_relate ⟶ wikistat_class
  {database = "", name = "wikistat_base", type = "fact_dimension_join", dimension_join = [
    [
      {data_source = "wikistat",        dimension = ["project"]},
      {data_source = "wikistat_relate", dimension = ["project"]},
    ],
    [
      {data_source = "wikistat_relate", dimension = ["class_id"]},
      {data_source = "wikistat_class",  dimension = ["class_id"]},
    ],
  ]},
]

metrics = [
  # Expose fact-table metrics through the joined source
  {data_source = "wikistat_base", type = "METRIC_AS", name = "hits",  value_type = "VALUE_INTEGER", dependency = ["wikistat.hits"]},
]

dimensions = [
  {data_source = "wikistat",        type = "DIMENSION_SINGLE", name = "project",    field_name = "project", value_type = "VALUE_STRING"},
  {data_source = "wikistat_relate", type = "DIMENSION_SINGLE", name = "project",    field_name = "project", value_type = "VALUE_STRING"},
  {data_source = "wikistat_relate", type = "DIMENSION_SINGLE", name = "class_id",   field_name = "class",   value_type = "VALUE_INTEGER"},
  {data_source = "wikistat_class",  type = "DIMENSION_SINGLE", name = "class_id",   field_name = "id",      value_type = "VALUE_INTEGER"},
  {data_source = "wikistat_class",  type = "DIMENSION_SINGLE", name = "class_name", field_name = "name",    value_type = "VALUE_STRING"},

  # Re-export through the joined source
  {data_source = "wikistat_base", type = "DIMENSION_MULTI", name = "project",    value_type = "VALUE_STRING",  dependency = ["wikistat.project", "wikistat_relate.project"]},
  {data_source = "wikistat_base", type = "DIMENSION_MULTI", name = "class_name", value_type = "VALUE_STRING",  dependency = ["wikistat_class.class_name"]},
]
```

### Query

```go
query := &types.Query{
    DataSetName: "wikistat_join",    // refers to set name
    Metrics:     []string{"hits"},
    Dimensions:  []string{"project", "class_name"},
}
```

### Generated SQL

```sql
SELECT
    wikistat.project         AS project,
    wikistat_class.name      AS class_name,
    1.0 * SUM(wikistat.hits) AS hits
FROM wikistat AS wikistat
JOIN wikistat_relate AS wikistat_relate
  ON wikistat.project = wikistat_relate.project
JOIN wikistat_class AS wikistat_class
  ON wikistat_relate.class = wikistat_class.id
GROUP BY wikistat.project, wikistat_class.name
```

---

## 5. Synchronous vs streaming (RunSync vs RunChan)

For most queries use `RunSync` — it returns the complete result slice.  
For very large result sets (millions of rows), use `RunChan` to process rows as they stream in, avoiding a large in-memory buffer.

### RunSync (default)

```go
result, err := manager.RunSync(query)
if err != nil {
    log.Fatal(err)
}
// All rows are in result.Source
for _, row := range result.Source {
    fmt.Println(row["date"], row["hits"])
}
```

### RunChan (streaming)

```go
result, err := manager.RunChan(query)
if err != nil {
    log.Fatal(err)
}
// result.Source is populated row-by-row as the channel is drained internally.
// Use it the same way after RunChan returns:
for _, row := range result.Source {
    fmt.Println(row["date"], row["hits"])
}
```

> **Note**: `RunChan` buffers rows internally the same way as `RunSync` after it returns.
> The real benefit is that rows are fetched from the DB incrementally — the server starts
> sending data before all results are ready, which reduces time-to-first-row on large scans.

### Choosing between them

| | `RunSync` | `RunChan` |
|---|---|---|
| API simplicity | ✅ simpler | same after return |
| Memory for large results | loads all rows | fetches incrementally |
| Use case | typical queries | millions of rows / streaming ETL |

---

## 6. Debug: inspect generated SQL without running it

Use `BuildSQL` to see the SQL that would be executed, without actually running the query.  
This is useful for debugging, auditing, or logging.

```go
sql, err := manager.BuildSQL(query)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated SQL:\n", sql)
```

You can also enable GORM debug logging on all connections:

```go
import "gorm.io/gorm/logger"

manager.SetLogger(logger.Default.LogMode(logger.Info))
```

Or create the client with `Debug: true` to log every SQL statement:

```go
cfg := &olapsql.Configuration{
    ClientsOption: map[string]*olapsql.DBOption{
        "clickhouse": {
            DSN:   "clickhouse://localhost:9000/default",
            Type:  types.DBTypeClickHouse,
            Debug: true,   // ← prints SQL to stdout
        },
    },
    // ...
}
```

---

## 7. SQLite quick-start (no ClickHouse needed)

The test suite uses SQLite so you can experiment locally without any external database.

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    olapsql "github.com/awatercolorpen/olap-sql"
    "github.com/awatercolorpen/olap-sql/api/types"
)

func main() {
    dbPath := filepath.Join(os.TempDir(), "demo.db")

    cfg := &olapsql.Configuration{
        ClientsOption: map[string]*olapsql.DBOption{
            "sqlite": {
                DSN:  dbPath,
                Type: types.DBTypeSQLite,
            },
        },
        DictionaryOption: &olapsql.Option{
            AdapterOption: olapsql.AdapterOption{
                Type: olapsql.FILEAdapter,
                Dsn:  "test/dictionary.sqlite.toml", // included in the repo
            },
        },
    }

    manager, err := olapsql.NewManager(cfg)
    if err != nil {
        log.Fatal(err)
    }

    query := &types.Query{
        DataSetName: "wikistat",
        TimeInterval: &types.TimeInterval{
            Name:  "date",
            Start: "2021-05-06",
            End:   "2021-05-08",
        },
        Metrics:    []string{"hits"},
        Dimensions: []string{"date"},
        Orders:     []*types.OrderBy{{Name: "date", Direction: types.OrderDirectionTypeAscending}},
    }

    sql, err := manager.BuildSQL(query)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("SQL:", sql)
}
```

See `test/dictionary.sqlite.toml` in the repository for a complete working TOML schema, and the `*_test.go` files for more query examples.
