package models

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"strings"
)

var DefaultOlapSqlModelDataSourceTableName = "olap_sql_model_data_sources"

type DataSource struct {
	ID          uint64               `gorm:"column:id;primaryKey"    json:"id,omitempty"`
	Type        types.DataSourceType `gorm:"column:type"             json:"type"`
	Name        string               `gorm:"column:name;unique"      json:"name"`
	Description string               `gorm:"column:description"      json:"description"`
	Metrics     []*Metric            `gorm:"foreignKey:DataSourceID" json:"metrics,omitempty"`
	Dimensions  []*Dimension         `gorm:"foreignKey:DataSourceID" json:"dimensions,omitempty"`
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
