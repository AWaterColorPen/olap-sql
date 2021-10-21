package dictionary

import (
    "github.com/awatercolorpen/olap-sql/api/types"
    "strings"
)

type DataSet struct {
    Name        string           `toml:"name"`
    Description string           `toml:"description"`
    Join        []*DataSetJoin   `toml:"join"`
    Merged      []*DataSetMerged `toml:"merged"`
}


type DataSetMerged struct {
    DataSource []string `toml:"secondary"`

}

type DataSetJoin struct {
    DataSource1 string    `toml:"data_source1"`
    DataSource2 string    `toml:"data_source2"`
    JoinOn      []*JoinOn `toml:"join_on"`
}

type JoinOn struct {
    Dimension1 string `toml:"dimension1"`
    Dimension2 string `toml:"dimension2"`
}

type DataSource struct {
    Type        types.DataSourceType `toml:"type"`
    Name        string               `toml:"name"`
    Description string               `toml:"description"`
}

func (d *DataSource) GetTableName() string {
    out := strings.Split(d.Name, ".")
    return out[len(out)-1]
}

func (d *DataSource) GetDatabaseName() string {
    if !strings.Contains(d.Name, ".") {
        return ""
    }
    out := strings.Split(d.Name, ".")
    return out[0]
}

type Dimension struct {
    DataSource   string          `toml:"data_source"`
    Name         string          `toml:"name"`
    FieldName    string          `toml:"field_name"`
    ValueType    types.ValueType `toml:"value_type"`
    Description  string          `toml:"description"`
}

type Metric struct {
    DataSource   string           `toml:"data_source"`
    Name         string           `toml:"name"`
    FieldName    string           `toml:"field_name"`
    Type         types.MetricType `toml:"type"`
    ValueType    types.ValueType  `toml:"value_type"`
    Composition  []string         `toml:"composition"`
    Description  string           `toml:"description"`
}
