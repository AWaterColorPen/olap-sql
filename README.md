# olap-sql

[![Go](https://github.com/AWaterColorPen/olap-sql/actions/workflows/go.yml/badge.svg)](https://github.com/AWaterColorPen/olap-sql/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/awatercolorpen/olap-sql.svg)](https://pkg.go.dev/github.com/awatercolorpen/olap-sql)

## Introduction

**olap-sql** is a Go library that turns high-level OLAP query definitions into adapted SQL for multiple database backends (ClickHouse, MySQL, PostgreSQL, SQLite). You describe *what* you want — metrics, dimensions, filters — and olap-sql figures out *how* to query it.

### How it works

```
Query (metrics + dimensions + filters)
        ↓
  Dictionary (schema/config)
        ↓
  Clause (backend-specific IR)
        ↓
  SQL string  ──►  Database  ──►  Result
```

---

## Quick Start

### 1. Install

```bash
go get github.com/awatercolorpen/olap-sql
```

### 2. Define the schema (TOML)

Create `olap-sql.toml` describing your data model:

```toml
sets = [
  {name = "wikistat", type = "clickhouse", data_source = "wikistat"},
]

sources = [
  {database = "", name = "wikistat", type = "fact"},
]

metrics = [
  {data_source = "wikistat", type = "METRIC_SUM",    name = "hits",     field_name = "hits", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_COUNT",  name = "count",    field_name = "*",    value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_DIVIDE", name = "hits_avg", value_type = "VALUE_FLOAT", dependency = ["wikistat.hits", "wikistat.count"]},
]

dimensions = [
  {data_source = "wikistat", type = "DIMENSION_SINGLE", name = "date", field_name = "date", value_type = "VALUE_STRING"},
]
```

### 3. Create a Manager

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    olapsql "github.com/awatercolorpen/olap-sql"
    "github.com/awatercolorpen/olap-sql/api/types"
)

func main() {
    cfg := &olapsql.Configuration{
        // Map each DB type to a connection option.
        ClientsOption: olapsql.ClientsOption{
            "clickhouse": {
                DSN:  "clickhouse://localhost:9000/default",
                Type: types.DBTypeClickHouse,
            },
        },
        // Point to your TOML schema file.
        DictionaryOption: &olapsql.Option{
            AdapterOption: olapsql.AdapterOption{Dsn: "olap-sql.toml"},
        },
    }

    manager, err := olapsql.NewManager(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // --- Build the query ---
    queryJSON := `{
      "data_set_name": "wikistat",
      "time_interval": {"name": "date", "start": "2021-05-06", "end": "2021-05-08"},
      "metrics":    ["hits", "hits_avg"],
      "dimensions": ["date"]
    }`

    query := &types.Query{}
    if err := json.Unmarshal([]byte(queryJSON), query); err != nil {
        log.Fatal(err)
    }

    // --- (Optional) Inspect the generated SQL ---
    sql, err := manager.BuildSQL(query)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated SQL:", sql)

    // --- Run the query ---
    result, err := manager.RunSync(query)
    if err != nil {
        log.Fatal(err)
    }

    out, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(out))
}
```

**Generated SQL** (ClickHouse):

```sql
SELECT
  wikistat.date AS date,
  SUM(wikistat.hits) AS hits,
  (1.0 * SUM(wikistat.hits)) / NULLIF(COUNT(*), 0) AS hits_avg
FROM wikistat AS wikistat
WHERE wikistat.date >= '2021-05-06'
  AND wikistat.date < '2021-05-08'
GROUP BY wikistat.date
```

**Result JSON**:

```json
{
  "dimensions": ["date", "hits", "hits_avg"],
  "source": [
    {"date": "2021-05-06T00:00:00Z", "hits": 147,  "hits_avg": 49},
    {"date": "2021-05-07T00:00:00Z", "hits": 7178, "hits_avg": 897.25}
  ]
}
```

---

## Common Patterns

### Add filters

```go
query := &types.Query{
    DataSetName: "wikistat",
    Metrics:     []string{"hits"},
    Filters: []*types.Filter{
        {
            OperatorType: types.FilterOperatorTypeLessEquals,
            Name:         "date",
            Value:        []any{"2021-05-06"},
        },
    },
}
```

Generated SQL:

```sql
SELECT SUM(wikistat.hits) AS hits
FROM wikistat AS wikistat
WHERE wikistat.date <= '2021-05-06'
```

### Stream large result sets

For large queries, use `RunChan` to receive rows one at a time instead of buffering everything in memory:

```go
result, err := manager.RunChan(query)
```

### Inspect SQL without executing

Use `BuildSQL` to preview the generated query (useful for debugging):

```go
sql, err := manager.BuildSQL(query)
fmt.Println(sql)
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Getting Started](./docs/getting-started.md) | Step-by-step guide to your first query |
| [Configuration](./docs/configuration.md) | Configure Manager, clients, and the OLAP dictionary |
| [Query](./docs/query.md) | Define metrics, dimensions, filters, orders, and limits |
| [Result](./docs/result.md) | Parse and work with query results |
| [Examples](./docs/examples.md) | Common usage scenarios (ClickHouse joins, time filters, concurrency) |
| [API Reference](./docs/api.md) | Full public API — Manager, Query, Filter, Result |
| [Architecture](./docs/architecture.md) | Internal design for contributors |
| [Contributing](./CONTRIBUTING.md) | How to contribute to olap-sql |

---

## Requirements

- **Go 1.22+** (uses range-over-integer syntax)
- Supported databases: ClickHouse, MySQL, PostgreSQL, SQLite

---

## License

See the [License File](./LICENSE).
