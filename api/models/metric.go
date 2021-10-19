package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/types"
)

var DefaultOlapSqlModelMetricTableName = "olap_sql_model_metrics"

type Metric struct {
	ID           uint64           `gorm:"column:id;primaryKey"                                          json:"id,omitempty"`
	Type         types.MetricType `gorm:"column:type"                                                   json:"type"`
	Name         string           `gorm:"column:name;index:idx_olap_sql_model_metrics,unique"           json:"name"`
	FieldName    string           `gorm:"column:field_name"                                             json:"field_name"`
	ValueType    types.ValueType  `gorm:"column:value_type"                                             json:"value_type"`
	Composition  *Composition     `gorm:"column:composition"                                            json:"composition"`
	DataSourceID uint64           `gorm:"column:data_source_id;index:idx_olap_sql_model_metrics,unique" json:"data_source_id"`
	Description  string           `gorm:"column:description"                                            json:"description"`
}

func (Metric) TableName() string {
	return DefaultOlapSqlModelMetricTableName
}

type Composition struct {
	MetricID []uint64 `json:"metric_id"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (c *Composition) Scan(value interface{}) error {
	return scan(value, c)
}

// Value return json value, implement driver.Valuer interface
func (c Composition) Value() (driver.Value, error) {
	return value(c)
}

func scan(value interface{}, to interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, to)
}

func value(v interface{}) (driver.Value, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}
