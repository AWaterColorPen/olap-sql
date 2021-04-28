package models

import (
	"time"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

var DefaultOlapSqlModelMetricTableName = "olap_sql_model_metrics"

type Metric struct {
	ID           uint64               `gorm:"column:id;primaryKey"                                          json:"id,omitempty"`
	CreatedAt    time.Time            `gorm:"column:created_at"                                             json:"created_at,omitempty"`
	UpdatedAt    time.Time            `gorm:"column:updated_at"                                             json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt       `gorm:"column:delete_at;index"                                        json:"-"`
	Name         string               `gorm:"column:name;index:idx_olap_sql_model_metrics,unique"           json:"name"`
	FieldName    string               `gorm:"column:name"                                                   json:"field_name"`
	DataSourceID uint64               `gorm:"column:data_source_id;index:idx_olap_sql_model_metrics,unique" json:"data_source_id"`
	Description  string               `gorm:"column:description"                                            json:"description"`
	Type         types.DataSourceType `gorm:"column:type"                                                   json:"type"`
}

func (Metric) TableName() string {
	return DefaultOlapSqlModelMetricTableName
}
