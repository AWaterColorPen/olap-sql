package olapsql

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type DataDictionaryOption struct {
	DBOption
	DB *gorm.DB `json:"-"`
}

func (d *DataDictionaryOption) NewDB() (*gorm.DB, error) {
	if d.DB != nil {
		return d.DB, nil
	}
	return d.DBOption.NewDB()
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
	switch v := item.(type) {
	case *models.DataSet:
		if err := d.isValidDataSetSchema(v.Schema); err != nil {
			return err
		}
		return d.db.Create(item).Error
	default:
		return d.db.Create(item).Error
	}
}

func (d *DataDictionary) Update(item interface{}) error {
	switch v := item.(type) {
	case *models.DataSet:
		if err := d.isValidDataSetSchema(v.Schema); err != nil {
			return err
		}
		return d.db.Updates(item).Error
	default:
		return d.db.Updates(item).Error
	}
}

func (d *DataDictionary) Delete(item interface{}, id uint64) error {
	return d.db.Delete(item, id).Error
}

func (d *DataDictionary) Translator(query *types.Query) (Translator, error) {
	t := &dataDictionaryTranslator{db: d.db}
	if err := d.db.Take(&t.set, "name = ?", query.DataSet).Error; err != nil {
		return nil, err
	}

	if t.set.Schema == nil {
		return nil, fmt.Errorf("schema is nil for data_set %v", query.DataSet)
	}

	id := t.set.Schema.DataSourceID()
	if err := d.db.Find(&t.sources, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	if err := d.db.Find(&t.metrics, "data_source_id IN ?", id).Error; err != nil {
		return nil, err
	}

	if err := d.db.Find(&t.dimensions, "data_source_id IN ?", id).Error; err != nil {
		return nil, err
	}

	return t, nil
}

func (d *DataDictionary) isValidDataSetSchema(schema *models.DataSetSchema) error {
	// TODO
	return nil
}

func NewDataDictionary(option *DataDictionaryOption) (*DataDictionary, error) {
	db, err := option.NewDB()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.DataSet{}, &models.DataSource{}, &models.Metric{}, &models.Dimension{})
	if err != nil {
		return nil, err
	}

	return &DataDictionary{db: option.DB}, nil
}
