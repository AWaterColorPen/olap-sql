package models

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

var DefaultOlapSqlModelDimensionTableName = "olap_sql_model_dimensions"

type Dimension struct {
	ID           uint64          `toml:"id"             json:"id,omitempty"`
	Name         string          `toml:"name"           json:"name"`
	FieldName    string          `toml:"field_name"     json:"field_name"`
	ValueType    types.ValueType `toml:"value_type"     json:"value_type"`
	DataSourceID uint64          `toml:"data_source_id" json:"data_source_id"`
	Description  string          `toml:"description"    json:"description"`
}

func (Dimension) TableName() string {
	return DefaultOlapSqlModelDimensionTableName
}
