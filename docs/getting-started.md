# Getting Started with olap-sql

This guide walks you through everything you need to start generating OLAP SQL with olap-sql — from installation to running your first query.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Core Concepts](#core-concepts)
4. [Step 1 — Define your schema](#step-1--define-your-schema)
5. [Step 2 — Create a Manager](#step-2--create-a-manager)
6. [Step 3 — Build and inspect SQL](#step-3--build-and-inspect-sql)
7. [Step 4 — Run a query](#step-4--run-a-query)
8. [Step 5 — Work with the result](#step-5--work-with-the-result)
9. [What's Next](#whats-next)

---

## Prerequisites

- **Go 1.22+** — olap-sql uses range-over-integer syntax (Go 1.22) and generics (Go 1.18+).
- A running database. The examples below use **SQLite** (zero config, no server required) so you can run them immediately. ClickHouse, MySQL, and PostgreSQL are all supported — swap out the `ClientsOption` when you are ready.

---

## Installation

```bash
go get github.com/awatercolorpen/olap-sql
```

---

## Core Concepts

Before writing any code it helps to know the four moving parts:

| Concept        | What it is                                                          |
|----------------|---------------------------------------------------------------------|
| **Schema**     | A TOML file describing your tables, metrics, and dimensions.        |
| **Manager**    | The main entry point. Holds DB clients + the schema dictionary.     |
| **Query**      | A struct (or JSON) you build to describe *what* you want.           |
| **Result**     | The structured response: a list of column names and row data.       |

The flow is:

```
Query  →  Dictionary (schema)  →  Clause (IR)  →  SQL  →  DB  →  Result
```

You never write SQL by hand. You describe what you want and olap-sql figures out the SQL.

---

## Step 1 — Define your schema

Create a file called `olap-sql.toml` in your project root.

```toml
# Each "set" is a named query context (= a business data set).
sets = [
  {name = "sales", type = "sqlite", data_source = "orders"},
]

# Each "source" maps to a real database table.
sources = [
  {database = "", name = "orders", type = "fact"},
]

# Metrics define how numeric columns are aggregated.
metrics = [
  {data_source = "orders", type = "METRIC_COUNT", name = "order_count",  field_name = "*",      value_type = "VALUE_INTEGER"},
  {data_source = "orders", type = "METRIC_SUM",   name = "revenue",      field_name = "amount", value_type = "VALUE_FLOAT"},
  {data_source = "orders", type = "METRIC_DIVIDE", name = "avg_order",
    value_type = "VALUE_FLOAT", dependency = ["orders.revenue", "orders.order_count"]},
]

# Dimensions define how rows are grouped and filtered.
dimensions = [
  {data_source = "orders", type = "DIMENSION_SINGLE", name = "date",   field_name = "date",   value_type = "VALUE_STRING"},
  {data_source = "orders", type = "DIMENSION_SINGLE", name = "region", field_name = "region", value_type = "VALUE_STRING"},
]
```

**Key rules:**

- Each `sets` entry needs a `name` (used in queries), a `type` (database driver), and a `data_source` (references a `sources` entry).
- `metrics` with `type = "METRIC_DIVIDE"` don't have a `field_name` — they reference two other metrics via `dependency`.
- `dimensions` name must be unique within a `data_source`.

---

## Step 2 — Create a Manager

```go
package main

import (
    "log"

    olapsql "github.com/awatercolorpen/olap-sql"
    "github.com/awatercolorpen/olap-sql/api/types"
)

func main() {
    cfg := &olapsql.Configuration{
        // ClientsOption: one entry per database connection.
        // The key must match the "type" field in your sets (e.g. "sqlite").
        ClientsOption: olapsql.ClientsOption{
            "sqlite": {
                DSN:  "file::memory:?cache=shared",
                Type: types.DBTypeSQLite,
            },
        },
        // DictionaryOption: points to your schema TOML.
        DictionaryOption: &olapsql.Option{
            AdapterOption: olapsql.AdapterOption{Dsn: "olap-sql.toml"},
        },
    }

    manager, err := olapsql.NewManager(cfg)
    if err != nil {
        log.Fatal(err)
    }

    _ = manager // use it in the next steps
}
```

> **Tip:** You can register multiple DB connections by adding more keys to `ClientsOption`.  
> For example `"clickhouse"` and `"mysql"` can coexist — each set in the schema will route to its own client.

---

## Step 3 — Build and inspect SQL

Before executing anything, it is useful to see what SQL olap-sql would generate.

```go
query := &types.Query{
    DataSetName: "sales",
    TimeInterval: &types.TimeInterval{
        Name:  "date",
        Start: "2024-01-01",
        End:   "2024-02-01",
    },
    Metrics:    []string{"order_count", "revenue", "avg_order"},
    Dimensions: []string{"date", "region"},
    Orders: []*types.OrderBy{
        {Name: "date", Direction: types.OrderDirectionAscending},
    },
}

sql, err := manager.BuildSQL(query)
if err != nil {
    log.Fatal(err)
}
fmt.Println(sql)
```

Expected output (SQLite):

```sql
SELECT
  orders.date    AS date,
  orders.region  AS region,
  COUNT(*)                                           AS order_count,
  SUM(orders.amount)                                 AS revenue,
  (1.0 * SUM(orders.amount)) / NULLIF(COUNT(*), 0)  AS avg_order
FROM orders AS orders
WHERE (orders.date >= '2024-01-01' AND orders.date < '2024-02-01')
GROUP BY orders.date, orders.region
ORDER BY orders.date ASC
```

`BuildSQL` never touches the database — it is safe to call at any time and is very useful for debugging.

---

## Step 4 — Run a query

Once you are happy with the SQL, run it:

```go
result, err := manager.RunSync(query)
if err != nil {
    log.Fatal(err)
}
```

For **large result sets** where you don't want to buffer everything in memory, use the channel-based variant:

```go
result, err := manager.RunChan(query)
if err != nil {
    log.Fatal(err)
}
```

Both methods return `*types.Result`.

---

## Step 5 — Work with the result

```go
import (
    "encoding/json"
    "fmt"
)

// result.Dimensions is a []string listing every column in order.
fmt.Println("columns:", result.Dimensions)

// result.Source is []map[string]any — one map per row.
for _, row := range result.Source {
    fmt.Printf("date=%v  region=%v  revenue=%v  avg_order=%v\n",
        row["date"], row["region"], row["revenue"], row["avg_order"])
}

// Or marshal the whole result to JSON:
out, _ := json.MarshalIndent(result, "", "  ")
fmt.Println(string(out))
```

Example JSON output:

```json
{
  "dimensions": ["date", "region", "order_count", "revenue", "avg_order"],
  "source": [
    {"date": "2024-01-05", "region": "north", "order_count": 12, "revenue": 4320.0, "avg_order": 360.0},
    {"date": "2024-01-05", "region": "south", "order_count":  7, "revenue": 1890.5, "avg_order": 270.07},
    {"date": "2024-01-12", "region": "north", "order_count": 19, "revenue": 7600.0, "avg_order": 400.0}
  ]
}
```

---

## What's Next

| Topic | Where to look |
|-------|---------------|
| Full configuration reference | [docs/configuration.md](./configuration.md) |
| All query options (filters, orders, limits, raw SQL) | [docs/query.md](./query.md) |
| Result format details | [docs/result.md](./result.md) |
| Common usage scenarios | [docs/examples.md](./examples.md) |
| Architecture for contributors | [docs/architecture.md](./architecture.md) |
| Contribution guide | [CONTRIBUTING.md](../CONTRIBUTING.md) |
