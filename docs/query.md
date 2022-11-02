# Query

- [Protocol](#protocol)
    - [Overview](#overview)
    - [Data Set Name](#data-set-name)
    - [Time Interval](#time-interval)
    - [Metrics](#metrics)
    - [Dimensions](#dimensions)
    - [Filters](#filters)
    - [Order By](#order-by)
    - [Limit](#limit)
    - [Sql](#sql)

## Protocol

### Overview

`Query` is a golang structure as OLAP input, defining at [github.com/awatercolorpen/olap-sql/api/types](../api/types/query.go).

```golang
type Query struct {
	DataSetName  string        `json:"data_set_name"`
	TimeInterval *TimeInterval `json:"time_interval"`
	Metrics      []string      `json:"metrics"`
	Dimensions   []string      `json:"dimensions"`
	Filters      []*Filter     `json:"filters"`
	Orders       []*OrderBy    `json:"orders"`
	Limit        *Limit        `json:"limit"`
	Sql          string        `json:"Sql"`
}
```

Supported json protocol now.

#### Example:

```json
{
  "data_set_name": "",
  "time_interval": {},
  "metrics": [],
  "dimensions": [],
  "filters": [],
  "orders": [],
  "limit": {},
  "sql": ""
}
```

### Data Set Name

`data_set_name` **(required)** is the name of one business, also the name of [sets](./configuration.md#sets).

#### Example:

```json
{
  "data_set_name": "wikistat"
}
```

### Time Interval

`time_interval` **(optional)** is a special structure for time interval filter conditional.

It will auto translate to one [filters](#filters).

```sql
WHERE (`name` >= `start` AND `name` < `end`)
```

1. `name` **(required)** is the [dimension name](./configuration.md#dimensions) for time interval.
2. `start` **(required)** is a string that is valid for golang `time.Time` type.
3. `end` **(required)** is a string that is valid for golang `time.Time` type.

#### Example:

```json
{
  "time_interval": {
    "name": "date",
    "start": "2021-05-06",
    "end": "2021-05-08"
  }
}
```

### Metrics

`metrics` **(optional)** is list of [metrics name](./configuration.md#metrics).

At least one of metrics or dimensions is required.

#### Example:

```json
{
  "metrics": [
    "hits",
    "hits_avg"
  ]
}
```

### Dimensions

`dimensions` **(optional)** is list of [dimensions name](./configuration.md#dimensions).

At least one of metrics or dimensions is required.

#### Example:

```json
{
  "dimensions": [
    "date"
  ]
}
```

### Filters

#### Example:

```json
{
}
```

### Order By

`time_interval` **(optional)** is a special structure for time interval filter conditional.

It will auto translate to one [filters](#filters).

```sql
WHERE (`name` >= `start` AND `name` < `end`)
```

1. `name` **(required)** is the [dimension name](./configuration.md#dimensions) for time interval.
2. `start` **(required)** is a string that is valid for golang `time.Time` type.
3. `end` **(required)** is a string that is valid for golang `time.Time` type.

#### Example:

```json
{
  "time_interval": {
    "name": "date",
    "start": "2021-05-06",
    "end": "2021-05-08"
  }
}
```