# olap-sql

[![Go](https://github.com/AWaterColorPen/olap-sql/actions/workflows/go.yml/badge.svg)](https://github.com/AWaterColorPen/olap-sql/actions/workflows/go.yml)

## Introduction

olap-sql is golang library for generating **adapted sql** by **olap query** with metrics, dimension and filter. 
Then get **formatted sql result** by queried metrics and dimension.

### Example

There is unprocessed olap data with table named `wikistat`.

| date       | time                | hits |
|------------|---------------------|------|
| 2021-05-07 | 2021-05-07 09:28:27 | 4783 |
| 2021-05-07 | 2021-05-07 09:33:59 | 1842 |
| 2021-05-07 | 2021-05-07 10:34:12 | 0    |
| 2021-05-06 | 2021-05-06 20:32:41 | 5    |
| 2021-05-06 | 2021-05-06 21:16:39 | 139  |

It wants a sql to query the data with `metrics: sum(hits) / count(*)` and `dimension: date`.

```sql
SELECT wikistat.date AS date, ( ( 1.0 * SUM(wikistat.hits) ) /  NULLIF(( COUNT(*) ), 0) ) AS hits_avg FROM wikistat AS wikistat GROUP BY wikistat.date
```

It wants a sql to query the data with `metrics: sum(hits)` and `filter: date <= '2021-05-06'`.

```sql
SELECT SUM(wikistat.hits) AS hits FROM wikistat AS wikistat WHERE wikistat.date <= '2021-05-06'
```

## Getting Started

### Define the OLAP dictionary configuration file

Create a new file for example named `olap-sql.toml` to define 
[sets](./docs/configuration.md#sets), 
[sources](./docs/configuration.md#sources), 
[metrics](./docs/configuration.md#metrics),
[dimensions](./docs/configuration.md#dimensions).

```toml
sets = [
  {name = "wikistat", type = "clickhouse", data_source = "wikistat"},
]

sources = [
  {database = "", name = "wikistat", type = "fact"},
]

metrics = [
  {data_source = "wikistat", type = "METRIC_SUM", name = "hits", field_name = "hits", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_COUNT", name = "count", field_name = "*", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_DIVIDE", name = "hits_avg", value_type = "VALUE_FLOAT", dependency = ["wikistat.hits", "wikistat.count"]},
]

dimensions = [
  {data_source = "wikistat", type = "DIMENSION_SINGLE", name = "date", field_name = "date", value_type = "VALUE_STRING"},
]
```

### To make use of olap-sql in golang

Create a new [manager instance](./docs/configuration.md#manager-configuration).

```golang
import "github.com/awatercolorpen/olap-sql"

// set clients option
clientsOption := map[string]*olapsql.DBOption{
	"clickhouse": &olapsql.DBOption{
		DSN:  "clickhouse://localhost:9000/default", 
		Type: "clickhouse"
	}
},

// set dictionary option
dictionaryOption := olapsql.AdapterOption{
	Dsn: "olap_sql.toml",
}

// build manager configuration
configuration := &olapsql.Configuration{
	ClientsOption:    clientsOption, 
	DictionaryOption: dictionaryOption,
}

// create a new manager instance
manager, err := olapsql.NewManager(configuration)
```

Build olap-sql [query](./docs/query.md).

```golang
import "github.com/awatercolorpen/olap-sql/api/types"

queryJson := `
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
}`

query := &types.Query{}
err := json.Unmarshal([]byte(queryJson), query)
```

Run query to get result from manager.

```golang
// run query with parallel chan
result, err := manager.RunChan(query)

// run query with sync
result, err := manager.RunSync(query)
```

### Generate SQL then format result

firstly, auto generate sql. [For detail](./docs/query.md#generate-sql-from-query).

Then get [result](./docs/result.md) json with `dimensions` property and `source` property.

```json
{
  "dimensions": [
    "date",
    "hits",
    "hits_avg"
  ],
  "source": [
    {
      "date": "2021-05-06T00:00:00Z",
      "hits": 147,
      "hits_avg": 49
    },
    {
      "date": "2021-05-07T00:00:00Z",
      "hits": 7178,
      "hits_avg": 897.25
    }
  ]
}
```

## Documentation

1. [Configuration](./docs/configuration.md) to configure olap-sql instance and OLAP dictionary.
2. [Query](./docs/query.md) to define olap query.
3. [Result](./docs/result.md) to parse olap result.

## License

See the [License File](./LICENSE).
