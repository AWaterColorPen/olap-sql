package olapsql_test

import (
	"fmt"
	"os"
	"time"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const mockWikiStatDataSet = "wikistat"

type WikiStat struct {
	Date       time.Time `gorm:"column:date"       json:"date"`
	Time       time.Time `gorm:"column:time"       json:"time"`
	Project    string    `gorm:"column:project"    json:"project"`
	SubProject string    `gorm:"column:subproject" json:"subproject"`
	Path       string    `gorm:"column:path"       json:"path"`
	Hits       uint64    `gorm:"column:hits"       json:"hits"`
	Size       uint64    `gorm:"column:size"       json:"size"`
}

func (WikiStat) TableName() string {
	return mockWikiStatDataSet
}

type WikiStatRelate struct {
	Project string  `gorm:"column:project" json:"project"`
	Class   uint64  `gorm:"column:class"   json:"class"`
	Source  float64 `gorm:"column:source"  json:"source"`
}

func (WikiStatRelate) TableName() string {
	return mockWikiStatDataSet + "_relate"
}

// curl -sSL "https://dumps.wikimedia.org/other/pagecounts-raw/2016/2016-01/" | grep -oE 'pagecounts-[0-9]+-[0-9]+\.gz' | sort | uniq | tee links.txt
func timeParseDate(in string) time.Time {
	t, _ := time.Parse("2006-01-02", in)
	return t
}

func timeParseTime(in string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", in)
	return t
}

func mockTimeGroupDimension(name, fieldName string, dataSourceID uint64) *models.Dimension {
	dimension := &models.Dimension{Name: name, ValueType: types.ValueTypeString, DataSourceID: dataSourceID}
	if DataWithClickhouse() {
		dimension.FieldName = fmt.Sprintf("formatDateTime(%v, '%%Y-%%m-%%d %%H:00:00')", fieldName)
		return dimension
	}

	dimension.FieldName = fmt.Sprintf("strftime('%%Y-%%m-%%d %%H', %v)", fieldName)
	return dimension
}

func MockWikiStatData(db *gorm.DB) error {
	if DataWithClickhouse() {
		return nil
	}

	if err := db.AutoMigrate(&WikiStat{}, &WikiStatRelate{}); err != nil {
		return err
	}
	if err := db.Debug().Create([]*WikiStat{
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T11:45:26Z"), Project: "city", SubProject: "CHN", Path: "level1", Hits: 121, Size: 4098},
		{Date: timeParseDate("2021-05-06"), Time: timeParseTime("2021-05-06T11:45:26Z"), Project: "city", SubProject: "CHN", Path: "level1", Hits: 139, Size: 10086},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T12:43:56Z"), Project: "city", SubProject: "CHN", Path: "level2", Hits: 20, Size: 1024},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T07:00:12Z"), Project: "city", SubProject: "US", Path: "level1", Hits: 19, Size: 2048},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T21:23:48Z"), Project: "school", SubProject: "university", Path: "engineering", Hits: 2, Size: 156},
		{Date: timeParseDate("2021-05-06"), Time: timeParseTime("2021-05-06T21:16:39Z"), Project: "school", SubProject: "university", Path: "engineering", Hits: 3, Size: 158},
		{Date: timeParseDate("2021-05-06"), Time: timeParseTime("2021-05-06T20:32:41Z"), Project: "school", SubProject: "senior", Path: "*", Hits: 5, Size: 212},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T09:28:27Z"), Project: "music", SubProject: "pop", Path: "", Hits: 4783, Size: 37291},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T09:31:23Z"), Project: "music", SubProject: "pop", Path: "ancient", Hits: 391, Size: 2531},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T09:33:59Z"), Project: "music", SubProject: "rap", Path: "", Hits: 1842, Size: 12942},
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T09:34:12Z"), Project: "music", SubProject: "rock", Path: "", Hits: 1093, Size: 9023},
	}).Error; err != nil {
		return err
	}
	if err := db.Debug().Create([]*WikiStatRelate{
		{Project: "city", Class: 1, Source: 4.872},
		{Project: "school", Class: 1, Source: 0.18742},
		{Project: "food", Class: 2, Source: 10.2484},
		{Project: "person", Class: 3, Source: 1.73},
		{Project: "music", Class: 4, Source: 93.20},
		{Project: "company", Class: 5, Source: 0.0281},
	}).Error; err != nil {
		return err
	}
	return nil
}

func MockWikiStatDataDictionary(dictionary *olapsql.DataDictionary) error {
	if err := dictionary.Create([]*models.DataSource{
		{Type: types.DataSourceTypeClickHouse, Name: mockWikiStatDataSet},
		{Type: types.DataSourceTypeClickHouse, Name: mockWikiStatDataSet + "_relate"},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.Metric{
		{Type: types.MetricTypeSum, Name: "hits_sum", FieldName: "hits", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{Type: types.MetricTypeSum, Name: "size_sum", FieldName: "size", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{Type: types.MetricTypeCount, Name: "count", FieldName: "*", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{Type: types.MetricTypeDivide, Name: "hits_avg", Composition: &models.Composition{MetricID: []uint64{1, 3}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{Type: types.MetricTypeDivide, Name: "size_avg", Composition: &models.Composition{MetricID: []uint64{2, 3}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{Type: types.MetricTypeDivide, Name: "hits_per_size", Composition: &models.Composition{MetricID: []uint64{1, 2}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{Type: types.MetricTypeSum, Name: "source_sum", FieldName: "source", ValueType: types.ValueTypeFloat, DataSourceID: 2},
		{Type: types.MetricTypeCount, Name: "count", FieldName: "*", ValueType: types.ValueTypeInteger, DataSourceID: 2},
		{Type: types.MetricTypeDivide, Name: "source_avg", Composition: &models.Composition{MetricID: []uint64{7, 8}}, ValueType: types.ValueTypeFloat, DataSourceID: 2},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.Dimension{
		{Name: "date", FieldName: "date", ValueType: types.ValueTypeString, DataSourceID: 1},
		mockTimeGroupDimension("time_by_hour", "time", 1),
		{Name: "project", FieldName: "project", ValueType: types.ValueTypeString, DataSourceID: 1},
		{Name: "sub_project", FieldName: "subproject", ValueType: types.ValueTypeString, DataSourceID: 1},
		{Name: "path", FieldName: "path", ValueType: types.ValueTypeString, DataSourceID: 1},
		{Name: "project", FieldName: "project", ValueType: types.ValueTypeString, DataSourceID: 2},
		{Name: "class", FieldName: "class", ValueType: types.ValueTypeInteger, DataSourceID: 2},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.DataSet{
		{
			Name: mockWikiStatDataSet,
			Schema: &models.DataSetSchema{PrimaryID: 1, Secondary: []*models.Secondary{
				{DataSourceID: 2, JoinOn: []*models.JoinOn{{DimensionID1: 3, DimensionID2: 6}}},
			}},
		},
	}); err != nil {
		return err
	}

	return nil
}

func MockLoad(manager *olapsql.Manager) error {
	dictionary, _ := manager.GetDataDictionary()
	if err := MockWikiStatDataDictionary(dictionary); err != nil {
		return err
	}

	client, _ := manager.GetClients()
	if db, err := client.Get(types.DataSourceTypeClickHouse, ""); err == nil {
		if err := MockWikiStatData(db); err != nil {
			return err
		}	
	}

	return nil
}

func MockQuery1() *types.Query {
	query := &types.Query{
		DataSetName: mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:    []string{"hits_sum", "size_sum", "hits_avg", "hits_per_size", "source_avg"},
		Dimensions: []string{"date", "class"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeNotIn, Name: "path", Value: []interface{}{"*"}},
			{OperatorType: types.FilterOperatorTypeIn, Name: "class", Value: []interface{}{1, 2, 3, 4}},
		},
	}
	return query
}

func MockQuery1ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 7)
	assert.Equal(t,"date", result.Dimensions[0])
	assert.Equal(t,"source_avg", result.Dimensions[6])
	assert.Len(t, result.Source, 3)
	assert.Len(t, result.Source[0], 7)
	assert.Equal(t, float64(10244), result.Source[0]["size_sum"])
	assert.Equal(t, 0.013861772745021476, result.Source[0]["hits_per_size"])
	assert.Equal(t, 2.52971, result.Source[0]["source_avg"])
}

func MockQuery2() *types.Query {
	query := &types.Query{
		DataSetName: mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:    []string{"hits_sum", "size_sum", "hits_avg", "hits_per_size", "source_avg"},
		Dimensions: []string{"time_by_hour", "class"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeNotIn, Name: "path", Value: []interface{}{"*"}},
			{OperatorType: types.FilterOperatorTypeIn, Name: "class", Value: []interface{}{1, 2, 3, 4}},
		},
	}
	return query
}

func MockQuery2ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 7)
	assert.Equal(t,"time_by_hour", result.Dimensions[0])
	assert.Equal(t,"source_avg", result.Dimensions[6])
	assert.Len(t, result.Source, 7)
	assert.Len(t, result.Source[0], 7)
	assert.Equal(t, float64(10086), result.Source[0]["size_sum"])
	assert.Equal(t, 0.013781479278207416, result.Source[0]["hits_per_size"])
	assert.Equal(t, 4.872, result.Source[0]["source_avg"])
}

func DataWithClickhouse() bool {
	args := os.Args
	for _, arg := range args {
		if arg == "clickhouse" {
			return true
		}
	}
	return false
}
