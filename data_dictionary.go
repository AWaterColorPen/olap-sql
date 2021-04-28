package olapsql

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

func (d *DataDictionaryOption) NewDB() (*gorm.DB, error) {
	if d.DB != nil {
		return d.DB, nil
	}

	dialect, err := getDialect(d.Type, d.DSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if d.Debug {
		db = db.Debug()
	}

	return db, nil
}

func getDialect(ty DBType, dsn string) (gorm.Dialector, error) {
	switch ty {
	case DBTypeSQLite:
		return sqlite.Open(dsn), nil
	case DBTypeMySQL:
		return mysql.Open(dsn), nil
	case DBTypePostgre:
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %v", ty)
	}
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
	db, err := option.NewDB()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.DataSet{}, &models.DataSource{}, &models.Metric{}, &models.Dimension{});
	if err != nil {
		return nil, err
	}

	return &DataDictionary{db: option.DB}, nil
}
