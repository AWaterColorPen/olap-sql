package models

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

var DefaultOlapSqlModelDataSetTableName = "olap_sql_model_data_sets"

type DataSet struct {
	ID          uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name        string         `gorm:"column:name;unique"     json:"name"`
	Description string         `gorm:"column:description"     json:"description"`
	Schema      *DataSetSchema `gorm:"column:schema"          json:"schema"`
}

func (DataSet) TableName() string {
	return DefaultOlapSqlModelDataSetTableName
}

type DataSetSchema struct {
	PrimaryID uint64       `json:"primary_id"`
	Secondary []*Secondary `json:"secondary"`
}

func (d *DataSetSchema) DataSourceID() []uint64 {
	id := []uint64{d.PrimaryID}
	for _, v := range d.Secondary {
		id = append(id, v.DataSourceID)
	}
	return id
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DataSetSchema) Scan(value interface{}) error {
	return scan(value, d)
}

// Value return json value, implement driver.Valuer interface
func (d DataSetSchema) Value() (driver.Value, error) {
	return value(d)
}

type Secondary struct {
	DataSourceID uint64    `json:"data_source_id"`
	JoinOn       []*JoinOn `json:"join_on"`
}

type JoinOn struct {
	DimensionID1 uint64 `json:"dimension_id1"`
	DimensionID2 uint64 `json:"dimension_id2"`
}
