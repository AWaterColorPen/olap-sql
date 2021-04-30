package olapsql_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewDataDictionary(t *testing.T) {
	d, err := newDataDictionary(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, dataDictionaryMockData(d))
}

func TestDataDictionary_Translator(t *testing.T) {
	d, err := newDataDictionary(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, dataDictionaryMockData(d))

	query := &types.Query{
		Metrics: []string{"cost", "click"},
		Dimensions: []string{"account_id", "account_name"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeIn, Name: "partition_time", Value: []interface{}{"2021-04-29", "2021-04-30"}},
		},
		DataSet: "account-cost",
	}

	translator, err := d.Translator(query)
	assert.NoError(t, err)
	request, err := translator.Translate(query)
	assert.NoError(t, err)

	sql, err := request.(*types.Request).Statement()
	assert.NoError(t, err)
	t.Log(sql)
}

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func newDataDictionary(sqlitePath string) (*olapsql.DataDictionary, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	option := &olapsql.DataDictionaryOption{DB: db}
	return olapsql.NewDataDictionary(option)
}

func dataDictionaryMockData(dictionary *olapsql.DataDictionary) error {
	if err := dictionary.Create([]*models.DataSource{
		{Type: types.DataSourceTypeClickHouse, Name: "mock-table"},
		{Type: types.DataSourceTypeClickHouse, Name: "mock-secondary-table-account"},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.Metric{
		{Type: types.MetricTypeSum, Name: "cost", FieldName: "cost", ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{Type: types.MetricTypeSum, Name: "click", FieldName: "click_count", ValueType: types.ValueTypeInteger, DataSourceID: 1},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.Dimension{
		{Name: "partition_time", FieldName: "partition", ValueType: types.ValueTypeString, DataSourceID: 1},
		{Name: "account_id", FieldName: "account_id", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{Name: "account_id", FieldName: "account_id", ValueType: types.ValueTypeInteger, DataSourceID: 2},
		{Name: "account_name", FieldName: "name", ValueType: types.ValueTypeString, DataSourceID: 2},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.DataSet{
		{
			Name: "account-cost",
			Schema: &models.DataSetSchema{PrimaryID: 1, Secondary: []*models.Secondary{
				{DataSourceID: 2, JoinOn: []*models.JoinOn{{DimensionID1: 2, DimensionID2: 3}}},
			}},
		},
	}); err != nil {
		return err
	}

	return nil
}
