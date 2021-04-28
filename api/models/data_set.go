package models

import (
	"gorm.io/gorm"
	"time"
)

var DefaultOlapSqlModelDataSetTableName = "olap_sql_model_data_sets"

type DataSet struct {
	ID        uint64               `gorm:"column:id;primaryKey"      json:"id,omitempty"`
	CreatedAt time.Time            `gorm:"column:created_at"         json:"created_at,omitempty"`
	UpdatedAt time.Time            `gorm:"column:updated_at"         json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt       `gorm:"column:delete_at;index" json:"-"`
	Name      string               `gorm:"column:name"            json:"name"`
}

func (DataSet) TableName() string {
	return DefaultOlapSqlModelDataSetTableName
}