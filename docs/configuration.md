# Configuration

- [Manager configuration](#manager-configuration)
  - [Database client](#database-client)
  - [Dictionary option](#dictionary-option)
  - [Create manager](#create-manager)
- [OLAP dictionary configuration](#olap-dictionary-configuration)
  - [Sets](#sets)
  - [Sources](#sources)
  - [Metrics](#metrics)
  - [Dimensions](#dimensions)

## Manager configuration

Configuration for creating a new manager instance.

### Database client

Configuration for defining olap database client.
It is a map of `key = name` and `value = one client` structure.
Each client have `dsn` and `type` two properties.

#### Supported OLAP database client type:

| type         | [dsn relate package version](../go.mod) |
|--------------|-----------------------------------------|
| `clickhouse` | gorm.io/driver/clickhouse v0.5.0        |
| `sqlite`     | gorm.io/driver/sqlite v1.3.6            |
| `mysql`      | gorm.io/driver/mysql v1.3.6             |
| `postgres`   | gorm.io/driver/postgres v1.3.6          |

#### Example:

```json
{
  "clients_option": {
    "clickhouse": {
      "dsn": "clickhouse://localhost:9000/default",
      "type": "clickhouse"
    }
  }
}
```

### Dictionary option

Option for loading [OLAP dictionary configuration](#olap-dictionary-configuration).
`dsn` is the path to load OLAP dictionary configuration.
`type` is dictionary adaptor type.

#### Supported dictionary option adaptor type:

| adaptor type | description                                  |
|--------------|----------------------------------------------|
| ` `          | default type is `FILE`                       |
| `FILE`       | load OLAP dictionary configuration from file |

#### Example:

```json
{
  "dictionary_option": {
    "dsn": "olap_sql.toml",
    "type": "FILE"
  }
}
```

### Create manager

Create a manager instance by  clients option and dictionary option.

```golang
import "github.com/awatercolorpen/olap-sql"

// set clients option
clientsOption := map[string]*olapsql.DBOption{},

// set dictionary option
dictionaryOption := olapsql.AdapterOption{}

// build manager configuration
configuration := &olapsql.Configuration{
	ClientsOption:    clientsOption, 
	DictionaryOption: dictionaryOption,
}

// create a new manager instance
manager, err := olapsql.NewManager(configuration)
if err != nil {
	log.Fatal(err)
}
```

## OLAP dictionary configuration

### Sets

Each set is one set of business query.

1. `name` is the name of business. It is used by [query's data_set_name](query.md#data-set-name) property.
2. `type` is the OLAP database client type.
3. `data_source` is the name of [sources](#sources).

```toml
sets = [
  {name = "wikistat", type = "clickhouse", data_source = "wikistat"},
]
```

#### Supported OLAP database client type:

Same with [database client](#database-client) type

### Sources

Each source defines one table or the relationship between tables.

1. `database` is the database name of source. The *dimension type* source must set `database` property.
2. `name` is the name of source. It is used by [sets](#sets) `data_source` property.
3. `type` is the source / table type.

```toml
sources = [
  {database = "", name = "wikistat", type = "fact"},
]
```

#### Supported source type:

| type                  | description                                           |
|-----------------------|-------------------------------------------------------|
| `dimension`           | dimension table, used by `JOIN ON`                    |
| `fact`                | fact table with really olap raw data                  |
| `fact_dimension_join` | combination table from fact table and dimension table |
| `merge_join`          | `JOIN ON` two different fact tables                   |

### Metrics

Each metrics defines how one metric be calculated. It is used for `SELECT` `WHERE` `ORDER BY`.

1. `data_source` is the name of [sources](#sources). This metrics belong to only one source.
2. `type` is the metrics type.
3. `name` is the name of metrics. It is used by query's `metrics` property.
4. `field_name` is column name in database table.
5. `value_type` is the value type. 
6. `dependency` is used for the *composition type* metrics.

```toml
metrics = [
  {data_source = "wikistat", type = "METRIC_SUM", name = "hits", field_name = "hits", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_COUNT", name = "count", field_name = "*", value_type = "VALUE_INTEGER"},
  {data_source = "wikistat", type = "METRIC_DIVIDE", name = "hits_avg", value_type = "VALUE_FLOAT", dependency = ["wikistat.hits", "wikistat.count"]},
]
```

#### Supported metric type:

| type                    | description                                         | classify         | required                                                                  | example                                |
|-------------------------|-----------------------------------------------------|------------------|---------------------------------------------------------------------------|----------------------------------------|
| `METRIC_VALUE`          | original value as one metrics.                      | single type      | must set `field_name`. eg `field_name = "value"`                          | `SELECT value AS name`                 |
| `METRIC_COUNT`          | count one column as one metrics.                    | single type      | must set `field_name`. eg `field_name = "*"`                              | `SELECT COUNT(*) AS name`              |
| `METRIC_DISTINCT_COUNT` | distinct count one metrics to a new metrics.        | single type      | must set `field_name`. eg `field_name = "value"`                          | `SELECT COUNT(DISTINCT value) AS name` |
| `METRIC_SUM`            | sum the metrics to a new metrics.                   | single type      | must set `field_name`. eg `field_name = "value"`                          | `SELECT SUM(value) AS name`            |
| `METRIC_ADD`            | add multi metrics to a new metrics.                 | composition type | must set `dependency`. eg `dependency = ["value1", "value2"]`             | `SELECT value1 + value2 AS name`       |
| `METRIC_SUBTRACT`       | subtract multi metrics to a new metrics.            | composition type | must set `dependency`. eg `dependency = ["value1", "value2"]`             | `SELECT value1 - value2 AS name`       |
| `METRIC_MULTIPLY`       | multiply multi metrics to a new metrics.            | composition type | must set `dependency`. eg `dependency = ["value1", "value2"]`             | `SELECT value1 * value2 AS name`       |
| `METRIC_DIVIDE`         | divide multi metrics to a new metrics.              | composition type | must set `dependency`. eg `dependency = ["value1", "value2"]`             | `SELECT value1 / value2 AS name`       |
| `METRIC_AS`             | relation with other table metrics to a new metrics. | composition type | must set `dependency`. eg `dependency = ["table.value1", "table.value2"]` | `SELECT table.value1 AS name`          |
| `METRIC_EXPRESSION`     | expression as one metrics                           | single type      | must set `field_name`. eg `field_name = "100"`                            | `SELECT 100 AS name`                   |

#### Supported metric value type:

Same with [database client](#database-client) type

### Dimensions

Each dimension defines OLAP data dimension to analyze. It is used for `SELECT` `JOIN ON` `WHERE` `GROUP BY` `ORDER BY`.

1. `data_source` is the name of [sources](#sources). This dimension belong to only one source.
2. `type` is the dimension type.
3. `name` is the name of dimension. It is used by query's `dimension` property.
4. `field_name` is column name in database table.
5. `value_type` is the value type.
6. `dependency` is used for the *composition type* metrics.

```toml
dimensions = [
  {data_source = "wikistat", type = "DIMENSION_SINGLE", name = "date", field_name = "date", value_type = "VALUE_STRING"},
]
```

#### Supported dimension type:

| type                   | description                                             | classify         | required                                                                    | example                                                               |
|------------------------|---------------------------------------------------------|------------------|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| `METRIC_VALUE`         | original value as one dimension.                        | single type      |                                                                             | `SELECT name ... GROUP BY name`                                       |
| `DIMENSION_SINGLE`     | one single column as one dimension.                     | single type      | must set `field_name`. eg `field_name = "value"`                            | `SELECT value AS name ... GROUP BY value`                             |
| `DIMENSION_MULTI`      | relation with other tables dimension to a new dimension | composition type | must set `dependency`. eg `dependency = ["table1.value"]`                   | `SELECT table1.value AS name ... GROUP BY table1.value`               |
| `DIMENSION_CASE`       | `CASE` pattern for multi dimension to a new dimension   | composition type | must set `dependency`. eg `dependency = ["table1.value1", "table2.value2"]` | `SELECT CASE WHEN table1.value1 != '' THEN table2.value2 END AS name` |
| `DIMENSION_EXPRESSION` | expression as one dimension.                            | single type      | must set `field_name`. eg `field_name = "formatDateTime(time, '%Y-%m-%d')"` | `SELECT formatDateTime(time, '%Y-%m-%d') AS name ... GROUP BY name`   |

#### Supported dimension value type:

Same with [database client](#database-client) type
