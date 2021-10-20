package models

import (
	"strings"

	"github.com/awatercolorpen/olap-sql/api/types"
)

var DefaultOlapSqlModelDataSourceTableName = "olap_sql_model_data_sources"

type DataSource struct {
	ID          uint64               `toml:"id"          json:"id,omitempty"`
	Type        types.DataSourceType `toml:"type"        json:"type"`
	Name        string               `toml:"name"        json:"name"`
	Description string               `toml:"description" json:"description"`
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

func (DataSource) TableName() string {
	return DefaultOlapSqlModelDataSourceTableName
}
