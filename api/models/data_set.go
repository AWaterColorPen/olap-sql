package models

import (
	"gorm.io/gorm"
	"time"
)

var DefaultOlapSqlModelDataSetTableName = "olap_sql_model_data_sets"

type DataSet struct {
	ID          uint64         `gorm:"column:id;primaryKey"                        json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"                           json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"                           json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index"                      json:"-"`
	Name        string         `gorm:"column:name;unique"                          json:"name"`
	Description string         `gorm:"column:description"                          json:"description"`
	PrimaryID   uint64         `gorm:"column:primary_id"                           json:"primary_id"`
	Primary     *DataSource    `gorm:"foreignKey:PrimaryID"                        json:"primary"`
	Secondary   []*DataSource  `gorm:"many2many:olap_sql_model_data_set_secondary" json:"secondary"`
}

func (DataSet) TableName() string {
	return DefaultOlapSqlModelDataSetTableName
}
