package models

import (
	"time"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

var DefaultOlapSqlModelDataSourceTableName = "olap_sql_model_data_sources"

type DataSource struct {
	ID          uint64               `gorm:"column:id;primaryKey"    json:"id,omitempty"`
	CreatedAt   time.Time            `gorm:"column:created_at"       json:"created_at,omitempty"`
	UpdatedAt   time.Time            `gorm:"column:updated_at"       json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt       `gorm:"column:delete_at;index"  json:"-"`
	Type        types.DataSourceType `gorm:"column:type"             json:"type"`
	Name        string               `gorm:"column:name;unique"      json:"name"`
	Description string               `gorm:"column:description"      json:"description"`
	Metrics     []*Metric            `gorm:"foreignKey:DataSourceID" json:"metrics,omitempty"`
	Dimension   []*Dimension         `gorm:"foreignKey:DataSourceID" json:"dimensions,omitempty"`
}

func (DataSource) TableName() string {
	return DefaultOlapSqlModelDataSourceTableName
}