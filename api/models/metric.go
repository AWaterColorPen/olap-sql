package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/types"
)

var DefaultOlapSqlModelMetricTableName = "olap_sql_model_metrics"

type Metric struct {
	ID           uint64           `yaml:"id"             json:"id,omitempty"`
	Type         types.MetricType `yaml:"type"           json:"type"`
	Name         string           `yaml:"name"           json:"name"`
	FieldName    string           `yaml:"field_name"     json:"field_name"`
	ValueType    types.ValueType  `yaml:"value_type"     json:"value_type"`
	Composition  *Composition     `yaml:"composition"    json:"composition"`
	DataSourceID uint64           `yaml:"data_source_id" json:"data_source_id"`
	Description  string           `yaml:"description"    json:"description"`
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
