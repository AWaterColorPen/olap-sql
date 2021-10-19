package models

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

var DefaultOlapSqlModelDimensionTableName = "olap_sql_model_dimensions"

type Dimension struct {
	ID           uint64          `yaml:"id"             json:"id,omitempty"`
	Name         string          `yaml:"name"           json:"name"`
	FieldName    string          `yaml:"field_name"     json:"field_name"`
	ValueType    types.ValueType `yaml:"value_type"     json:"value_type"`
	DataSourceID uint64          `yaml:"data_source_id" json:"data_source_id"`
	Description  string          `yaml:"description"    json:"description"`
}

func (Dimension) TableName() string {
	return DefaultOlapSqlModelDimensionTableName
}
