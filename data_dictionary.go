package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/models"
	"gorm.io/gorm"
)

type DBType string

const (
	DBTypeSQLite  DBType = "sqlite"
	DBTypeMySQL   DBType = "mysql"
	DBTypePostgre DBType = "postgres"
)

type DataDictionaryOption struct {
	Debug bool     `json:"debug"`
	DSN   string   `json:"dsn"`
	Type  DBType   `json:"type"`
	DB    *gorm.DB `json:"-"`
}

type DataDictionary struct {
	db *gorm.DB
}

func (d *DataDictionary) GetMetrics() ([]*models.Metric, error) {
	var out []*models.Metric
	return out, d.db.Find(&out).Error
}

func (d *DataDictionary) GetDimensions() ([]*models.Dimension, error) {
	var out []*models.Dimension
	return out, d.db.Find(&out).Error
}

func (d *DataDictionary) GetDataSources() ([]*models.DataSource, error) {
	var out []*models.DataSource
	return out, d.db.Find(&out).Error
}

func (d *DataDictionary) GetDataSets() ([]*models.DataSet, error) {
	var out []*models.DataSet
	return out, d.db.Find(&out).Error
}

func (d *DataDictionary) Create(item interface{}) error {
	return d.db.Create(item).Error
}

func (d *DataDictionary) Update(item interface{}) error {
	return d.db.Updates(item).Error
}

func (d *DataDictionary) Delete(item interface{}, id uint64) error {
	return d.db.Delete(item, id).Error
}

func NewDataDictionary(option *DataDictionaryOption) (*DataDictionary, error) {
	db := option.DB
	if db == nil {
		return nil, nil
	}

	if err := db.AutoMigrate(&models.DataSet{}, &models.DataSource{}, &models.Metric{}, &models.Dimension{}); err != nil {
		return nil, err
	}

	return &DataDictionary{db: option.DB}, nil
}
