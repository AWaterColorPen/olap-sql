package models

import (
	"time"

	"gorm.io/gorm"
)

var DefaultOlapSqlModelDimensionTableName = "olap_sql_model_dimensions"

type Dimension struct {
	ID           uint64         `gorm:"column:id;primaryKey"                                             json:"id,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at"                                                json:"created_at,omitempty"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"                                                json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"column:delete_at;index"                                           json:"-"`
	Name         string         `gorm:"column:name;index:idx_olap_sql_model_dimensions,unique"           json:"name"`
	FieldName    string         `gorm:"column:field_name"                                                json:"field_name"`
	DataSourceID uint64         `gorm:"column:data_source_id;index:idx_olap_sql_model_dimensions,unique" json:"data_source_id"`
	Description  string         `gorm:"column:description"                                               json:"description"`
}

func (Dimension) TableName() string {
	return DefaultOlapSqlModelDimensionTableName
}
