# Result

- [Protocol](#protocol)
    - [Overview](#overview)
    - [Dimensions](#dimensions)
    - [Source](#source)

## Protocol

### Overview

`Result` is a golang structure as OLAP output, defining at [github.com/awatercolorpen/olap-sql/api/types](../api/types/result.go).

```golang
type Result struct {
	Dimensions []string         `json:"dimensions"`
	Source     []map[string]any `json:"source"`
}
```

Supported json protocol now.

#### Example:

```json
{
  "dimensions": [],
  "source": []
}
```

### Dimensions

`dimensions` equals [`query.dimensions`](./query.md#dimensions) and [`query.metrics`](./query.md#metrics)

#### Example:

```json
{
  "dimensions": [
    "date",
    "hits",
    "hits_avg"
  ]
}
```

### Source

`Source` is an array of `map[string]any`. Each item is one row result from OLAP database query. 
The key is column name, the value is column value.

#### Example:

```json
{
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