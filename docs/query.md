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
- [Generate SQL from query](#generate-sql-from-query)

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

`time_interval` **(optional)** is a special structure for time interval filter condition.

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

`metrics` **(optional)** is a list of [metrics name](./configuration.md#metrics).

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

`dimensions` **(optional)** is a list of [dimensions name](./configuration.md#dimensions).

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

`filters` **(optional)** is a list of `WHERE` statement object.

1. `table` **(don't fill)** is an autofilled property by `name`.
2. `name` **(required)** is a valid [dimension name](./configuration.md#dimensions) or [metrics name](./configuration.md#metrics).
3. `field_property` **(don't fill)** is an autofilled property by `name`.
4. `operator_type` **(required)** is operator type of filter. [For detail](#supported-filter-operator-type-and-value).
5. `value` **(required)** is array for values. [For detail](#supported-filter-operator-type-and-value).
6. `value_type` **(don't fill)** is an autofilled property by `name`.
7. `children` **(optional)** is a list of [Filter](#filters). It is used for `FILTER_OPERATOR_AND` or `FILTER_OPERATOR_OR` case.

#### Example:

```json
{
}
```

#### Supported filter operator type and value

| type                             | description                  | required                       | sql example                     |
|----------------------------------|------------------------------|--------------------------------|---------------------------------|
| `FILTER_OPERATOR_EQUALS`         | `=` condition.               | `value` size equals 1.         | `name = value[0]`               |
| `FILTER_OPERATOR_IN`             | `IN` condition.              | `value` size must larger 0.    | `name IN (value)`               |
| `FILTER_OPERATOR_NOT_IN`         | `NOT IN` condition.          | `value` size must larger 0.    | `name NOT IN (value)`           |
| `FILTER_OPERATOR_LESS_EQUALS`    | `<=` condition.              | `value` size equals 1.         | `name <= value[0]`              |
| `FILTER_OPERATOR_LESS`           | `<` condition.               | `value` size equals 1.         | `name < value[0]`               |
| `FILTER_OPERATOR_GREATER_EQUALS` | `>=` condition.              | `value` size equals 1.         | `name >= value[0]`              |
| `FILTER_OPERATOR_GREATER`        | `>` condition.               | `value` size equals 1.         | `name > value[0]`               |
| `FILTER_OPERATOR_LIKE`           | `like` condition.            | `value` size equals 1.         | `name like value[0]`            |
| `FILTER_OPERATOR_HAS`            | `has` condition.             | `value` size equals 1.         | `hash(name, value[0])`          |
| `FILTER_OPERATOR_EXTENSION`      | expression as one condition  | `value` size equals 1.         | `value[0]`                      |
| `FILTER_OPERATOR_AND`            | `AND` multi children filters | `children` size must larger 0. | `(children[0] AND children[1])` |
| `FILTER_OPERATOR_OR`             | `OR` multi children filters  | `children` size must larger 0. | `(children[0] OR children[1])`  |

#### Supported value type

| type            | description |
|-----------------|-------------|
| `VALUE_STRING`  | string      |
| `VALUE_INTEGER` | int64       |
| `VALUE_FLOAT`   | float       |

### Order By

`orders` **(optional)** is list of `ORDER BY` statement object.

1. `table` **(don't fill)** is an autofilled property by `name`.
2. `name` **(required)** is a valid [dimension name](./configuration.md#dimensions) or [metrics name](./configuration.md#metrics).
3. `field_property` **(don't fill)** is an autofilled property by `name`.
4. `direction` **(optional)** is direction type for `ORDER BY`.

#### Example:

```json
{
  "orders": [
    {
      "name": "date",
      "direction": "ORDER_DIRECTION_DESCENDING"
    }
  ]
}
```

#### Supported direction type:

| type                         | description |
|------------------------------|-------------|
| ` `                          | `ASC`       |
| `ORDER_DIRECTION_ASCENDING`  | `ASC`       |
| `ORDER_DIRECTION_DESCENDING` | `DESC`      |

### Limit

`limit` **(optional)** is a structure for setting `LIMIT` and `OFFSET`.

```sql
LIMIT 100 OFFSET 20
```

1. `limit` **(optional)** is uint64 for setting `LIMIT`.
2. `offset` **(optional)** is uint64 for setting `OFFSET`.

#### Example:

```json
{
  "limit": {
    "limit": 100,
    "offset": 20
  }
}
```

### SQL

`sql` **(optional)** is a special property to override other `metrics`, `dimensions`, `filters`, `orders` and `limit` properties.

It will not generate SQL.

#### Example:

```json
{
  "sql": "SELECT VERSION()"
}
```

## Generate SQL from query

#### Query Example:

```json
{
  "data_set_name": "wikistat",
  "time_interval": {
    "name": "date",
    "start": "2021-05-06",
    "end": "2021-05-08"
  },
  "metrics": [
    "hits",
    "hits_avg"
  ],
  "dimensions": [
    "date"
  ]
}
```

#### Auto SQL Example:

```sql
SELECT
    wikistat.date AS date,
    1.0 * SUM(wikistat.hits) AS hits,
    ( ( 1.0 * SUM(wikistat.hits) ) /  NULLIF(( COUNT(*) ), 0) ) AS hits_avg
FROM wikistat AS wikistat
WHERE
    ( wikistat.date >= '2021-05-06' AND wikistat.date < '2021-05-08' )
GROUP BY
    wikistat.date
```
