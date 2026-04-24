# API Reference

This page documents the public API of olap-sql. Each section covers a type or function with its parameters, return values, and typical usage.

---

## Table of Contents

- [NewManager](#newmanager)
- [Manager](#manager)
  - [RunSync](#runsync)
  - [RunChan](#runchan)
  - [BuildSQL](#buildsql)
  - [BuildTransaction](#buildtransaction)
  - [SetLogger](#setlogger)
  - [GetClients](#getclients)
  - [GetDictionary](#getdictionary)
- [Configuration](#configuration)
- [Query](#query)
  - [TimeInterval](#timeinterval)
  - [Filter](#filter)
    - [FilterOperatorType](#filteroperatortype)
    - [ValueType](#valuetype)
  - [OrderBy](#orderby)
    - [OrderDirectionType](#orderdirectiontype)
  - [Limit](#limit)
- [Result](#result)

---

## NewManager

```go
func NewManager(cfg *Configuration) (*Manager, error)
```

Creates and initialises a `Manager` from the provided `Configuration`.

| Parameter | Type             | Description                                    |
|-----------|------------------|------------------------------------------------|
| `cfg`     | `*Configuration` | Holds client DSNs and a dictionary option.     |

**Returns** a ready-to-use `*Manager`, or an error if any client DSN is invalid or the dictionary file cannot be parsed.

**At least one** of `ClientsOption` or `DictionaryOption` must be non-nil; if both are nil, the manager is created but no queries can be executed.

**Example:**

```go
cfg := &olapsql.Configuration{
    ClientsOption: olapsql.ClientsOption{
        "clickhouse": {
            DSN:  "clickhouse://localhost:9000/default",
            Type: types.DBTypeClickHouse,
        },
    },
    DictionaryOption: &olapsql.Option{
        AdapterOption: olapsql.AdapterOption{Dsn: "olap-sql.toml"},
    },
}
manager, err := olapsql.NewManager(cfg)
if err != nil {
    log.Fatal(err)
}
```

---

## Manager

`Manager` is the main entry point. It holds a set of registered database clients and an OLAP dictionary.

### RunSync

```go
func (m *Manager) RunSync(query *types.Query) (*types.Result, error)
```

Executes the query **synchronously** and returns the full result.

| Parameter | Type          | Description                    |
|-----------|---------------|--------------------------------|
| `query`   | `*types.Query`| The OLAP query to execute.     |

**Returns** a `*types.Result` containing column names (`Dimensions`) and row data (`Source`), or an error.

Use this for typical queries where the result set fits in memory.

---

### RunChan

```go
func (m *Manager) RunChan(query *types.Query) (*types.Result, error)
```

Executes the query and **streams rows** internally over a channel before assembling the result.

| Parameter | Type          | Description                    |
|-----------|---------------|--------------------------------|
| `query`   | `*types.Query`| The OLAP query to execute.     |

**Returns** a `*types.Result` (same structure as `RunSync`), or an error.

Prefer `RunChan` for large result sets to avoid peak memory pressure.

---

### BuildSQL

```go
func (m *Manager) BuildSQL(query *types.Query) (string, error)
```

Translates the query into its **SQL string** without executing it.

| Parameter | Type          | Description                    |
|-----------|---------------|--------------------------------|
| `query`   | `*types.Query`| The OLAP query to translate.   |

**Returns** the generated SQL string, or an error if translation fails.

Useful for debugging, audit logging, or displaying the query to end users.

**Example:**

```go
sql, err := manager.BuildSQL(query)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated SQL:", sql)
```

---

### BuildTransaction

```go
func (m *Manager) BuildTransaction(query *types.Query) (*gorm.DB, error)
```

Translates the query into a `*gorm.DB` ready to execute.

| Parameter | Type          | Description                    |
|-----------|---------------|--------------------------------|
| `query`   | `*types.Query`| The OLAP query to translate.   |

**Returns** a configured `*gorm.DB`, or an error.

Use this when you need direct access to the GORM object — for example to attach custom hooks, inspect the SQL via `ToSQL`, or integrate with an existing GORM session.

---

### SetLogger

```go
func (m *Manager) SetLogger(log logger.Interface)
```

Attaches a custom GORM logger to all registered database clients.

| Parameter | Type               | Description                             |
|-----------|--------------------|-----------------------------------------|
| `log`     | `logger.Interface` | A GORM-compatible logger implementation.|

Call this after `NewManager` to enable query logging, debug output, or custom log routing.

---

### GetClients

```go
func (m *Manager) GetClients() (Clients, error)
```

Returns the registered database clients. Returns an error if the manager has no `ClientsOption`.

---

### GetDictionary

```go
func (m *Manager) GetDictionary() (*Dictionary, error)
```

Returns the OLAP dictionary. Returns an error if the manager has no `DictionaryOption`.

---

## Configuration

```go
type Configuration struct {
    ClientsOption    ClientsOption
    DictionaryOption *Option
}
```

| Field              | Type             | Description                                                        |
|--------------------|------------------|--------------------------------------------------------------------|
| `ClientsOption`    | `ClientsOption`  | Map of database name → `*DBOption`. Each entry registers one DB.  |
| `DictionaryOption` | `*Option`        | Points to the TOML schema file via `AdapterOption.Dsn`.           |

### ClientsOption / DBOption

```go
type ClientsOption map[string]*DBOption

type DBOption struct {
    DSN  string          // e.g. "clickhouse://localhost:9000/default"
    Type types.DBType    // e.g. types.DBTypeClickHouse
}
```

Supported `DBType` values:

| Constant                  | Database    |
|---------------------------|-------------|
| `types.DBTypeClickHouse`  | ClickHouse  |
| `types.DBTypeMySQL`       | MySQL       |
| `types.DBTypePostgreSQL`  | PostgreSQL  |
| `types.DBTypeSQLite`      | SQLite      |

---

## Query

```go
type Query struct {
    DataSetName  string         `json:"data_set_name"`
    TimeInterval *TimeInterval  `json:"time_interval"`
    Metrics      []string       `json:"metrics"`
    Dimensions   []string       `json:"dimensions"`
    Filters      []*Filter      `json:"filters"`
    Orders       []*OrderBy     `json:"orders"`
    Limit        *Limit         `json:"limit"`
    Sql          string         `json:"Sql"`
}
```

| Field          | Type             | Required | Description                                                     |
|----------------|------------------|----------|-----------------------------------------------------------------|
| `DataSetName`  | `string`         | ✅       | Must match a `sets[].name` entry in your TOML schema.           |
| `TimeInterval` | `*TimeInterval`  | ❌       | Shorthand for a start/end filter on a time dimension.           |
| `Metrics`      | `[]string`       | ❌       | Metric names defined in the TOML schema.                        |
| `Dimensions`   | `[]string`       | ❌       | Dimension names defined in the TOML schema.                     |
| `Filters`      | `[]*Filter`      | ❌       | Arbitrary filter conditions (supports nesting via AND/OR).      |
| `Orders`       | `[]*OrderBy`     | ❌       | Sort order for the result set.                                  |
| `Limit`        | `*Limit`         | ❌       | Pagination — row limit and offset.                              |
| `Sql`          | `string`         | ❌       | Reserved field; set by the library during translation.          |

### TimeInterval

```go
type TimeInterval struct {
    Name  string `json:"name"`
    Start string `json:"start"`
    End   string `json:"end"`
}
```

Convenience wrapper that expands to two `Filter` entries (`>= Start` and `< End`) on the named dimension.

| Field   | Description                                         |
|---------|-----------------------------------------------------|
| `Name`  | Dimension name to filter on (e.g. `"date"`).        |
| `Start` | Inclusive lower bound (e.g. `"2021-05-06"`).        |
| `End`   | **Exclusive** upper bound (e.g. `"2021-05-08"`).    |

**Example — equivalent SQL fragment:**

```sql
WHERE date >= '2021-05-06' AND date < '2021-05-08'
```

---

### Filter

```go
type Filter struct {
    OperatorType  FilterOperatorType `json:"operator_type"`
    ValueType     ValueType          `json:"value_type"`
    Table         string             `json:"table"`
    Name          string             `json:"name"`
    FieldProperty FieldProperty      `json:"field_property"`
    Value         []any              `json:"value"`
    Children      []*Filter          `json:"children"`
}
```

| Field          | Description                                                                 |
|----------------|-----------------------------------------------------------------------------|
| `OperatorType` | Comparison operator (see `FilterOperatorType`).                             |
| `ValueType`    | How to quote values in SQL (see `ValueType`).                               |
| `Name`         | Dimension or metric name to filter on.                                      |
| `Value`        | One or more comparison values (slice for `IN`/`NOT IN`, single for others). |
| `Children`     | Nested filters; only used with `FILTER_OPERATOR_AND` / `FILTER_OPERATOR_OR`.|

#### FilterOperatorType

| Constant                         | SQL equivalent               |
|----------------------------------|------------------------------|
| `FILTER_OPERATOR_EQUALS`         | `field = value`              |
| `FILTER_OPERATOR_IN`             | `field IN (v1, v2, ...)`     |
| `FILTER_OPERATOR_NOT_IN`         | `field NOT IN (...)`         |
| `FILTER_OPERATOR_LESS_EQUALS`    | `field <= value`             |
| `FILTER_OPERATOR_LESS`           | `field < value`              |
| `FILTER_OPERATOR_GREATER_EQUALS` | `field >= value`             |
| `FILTER_OPERATOR_GREATER`        | `field > value`              |
| `FILTER_OPERATOR_LIKE`           | `field LIKE value`           |
| `FILTER_OPERATOR_HAS`            | `has(field, value)` (CK)     |
| `FILTER_OPERATOR_EXTENSION`      | raw SQL expression           |
| `FILTER_OPERATOR_AND`            | `( child1 AND child2 ... )`  |
| `FILTER_OPERATOR_OR`             | `( child1 OR child2 ... )`   |

#### ValueType

| Constant        | SQL quoting              |
|-----------------|--------------------------|
| `VALUE_STRING`  | `'value'` (quoted)       |
| `VALUE_INTEGER` | `value` (unquoted)       |
| `VALUE_FLOAT`   | `value` (unquoted)       |
| `VALUE_UNKNOWN` | auto-detect from Go type |

**Nested filter example (AND + OR):**

```go
filter := &types.Filter{
    OperatorType: types.FilterOperatorTypeAnd,
    Children: []*types.Filter{
        {
            OperatorType: types.FilterOperatorTypeEquals,
            Name:         "project",
            ValueType:    types.ValueTypeString,
            Value:        []any{"en"},
        },
        {
            OperatorType: types.FilterOperatorTypeGreater,
            Name:         "hits",
            ValueType:    types.ValueTypeInteger,
            Value:        []any{1000},
        },
    },
}
```

Generated SQL:
```sql
( project = 'en' AND hits > 1000 )
```

---

### OrderBy

```go
type OrderBy struct {
    Name      string             `json:"name"`
    Direction OrderDirectionType `json:"direction"`
}
```

| Field       | Description                                           |
|-------------|-------------------------------------------------------|
| `Name`      | Metric or dimension name to sort by.                  |
| `Direction` | `ORDER_DIRECTION_ASCENDING` or `ORDER_DIRECTION_DESCENDING`. |

#### OrderDirectionType

| Constant                      | SQL          |
|-------------------------------|--------------|
| `ORDER_DIRECTION_ASCENDING`   | `name ASC`   |
| `ORDER_DIRECTION_DESCENDING`  | `name DESC`  |

---

### Limit

```go
type Limit struct {
    Limit  uint64 `json:"limit"`
    Offset uint64 `json:"offset"`
}
```

| Field    | Description                                     |
|----------|-------------------------------------------------|
| `Limit`  | Maximum number of rows to return.               |
| `Offset` | Number of rows to skip (for pagination).        |

---

## Result

```go
type Result struct {
    Dimensions []string         `json:"dimensions"`
    Source     []map[string]any `json:"source"`
}
```

| Field        | Description                                                                                        |
|--------------|----------------------------------------------------------------------------------------------------|
| `Dimensions` | Ordered list of column names — dimensions first, then metrics (mirrors the query field order).    |
| `Source`     | Slice of row maps. Each map is `column_name → value`. Value type matches the database driver type.|

**Example response:**

```json
{
  "dimensions": ["date", "hits", "hits_avg"],
  "source": [
    {"date": "2021-05-06T00:00:00Z", "hits": 147,  "hits_avg": 49},
    {"date": "2021-05-07T00:00:00Z", "hits": 7178, "hits_avg": 897.25}
  ]
}
```
