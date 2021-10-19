package olapsql_test

import (
	"fmt"
	olapsql "github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/awatercolorpen/olap-sql/dictionary"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"os"
	"time"
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

type ClassRelate struct {
	ID   uint64 `gorm:"column:id"   json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (ClassRelate) TableName() string {
	return mockWikiStatDataSet + "_class"
}

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

	dimension.FieldName = fmt.Sprintf("strftime('%%Y-%%m-%%d %%H:00:00', %v)", fieldName)
	return dimension
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

func DataSourceType() types.DataSourceType {
	if DataWithClickhouse() {
		return types.DataSourceTypeClickHouse
	}
	return types.DataSourceTypeUnknown
}

func MockWikiStatData(db *gorm.DB) error {
	if DataWithClickhouse() {
		return nil
	}

	if err := db.Debug().AutoMigrate(&WikiStat{}, &WikiStatRelate{}, &ClassRelate{}); err != nil {
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
		{Date: timeParseDate("2021-05-07"), Time: timeParseTime("2021-05-07T10:34:12Z"), Project: "music", SubProject: "rock", Path: "", Hits: 0, Size: 0},
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
	if err := db.Debug().Create([]*ClassRelate{
		{ID: 1, Name: "location"},
		{ID: 2, Name: "life"},
		{ID: 3, Name: "culture"},
		{ID: 4, Name: "entertainment"},
		{ID: 5, Name: "social"},
	}).Error; err != nil {
		return err
	}
	return nil
}

func MockWikiStatDataDictionary(dictionary *dictionary.Dictionary) error {
	if err := dictionary.Create([]*models.DataSource{
		{ID: 1, Type: DataSourceType(), Name: mockWikiStatDataSet},
		{ID: 2, Type: DataSourceType(), Name: mockWikiStatDataSet + "_relate"},
		{ID: 3, Type: DataSourceType(), Name: mockWikiStatDataSet + "_class"},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.Metric{
		{ID:1, Type: types.MetricTypeSum, Name: "hits", FieldName: "hits", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{ID:2, Type: types.MetricTypeSum, Name: "size_sum", FieldName: "size", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{ID:3, Type: types.MetricTypeCount, Name: "count", FieldName: "*", ValueType: types.ValueTypeInteger, DataSourceID: 1},
		{ID:4, Type: types.MetricTypeDivide, Name: "hits_avg", Composition: &models.Composition{MetricID: []uint64{1, 3}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{ID:5, Type: types.MetricTypeDivide, Name: "size_avg", Composition: &models.Composition{MetricID: []uint64{2, 3}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{ID:6, Type: types.MetricTypeDivide, Name: "hits_per_size", Composition: &models.Composition{MetricID: []uint64{1, 2}}, ValueType: types.ValueTypeFloat, DataSourceID: 1},
		{ID:7,Type: types.MetricTypeSum, Name: "source_sum", FieldName: "source", ValueType: types.ValueTypeFloat, DataSourceID: 2},
		{ID:8,Type: types.MetricTypeCount, Name: "count", FieldName: "*", ValueType: types.ValueTypeInteger, DataSourceID: 2},
		{ID:9,Type: types.MetricTypeDivide, Name: "source_avg", Composition: &models.Composition{MetricID: []uint64{7, 8}}, ValueType: types.ValueTypeFloat, DataSourceID: 2},
	}); err != nil {
		return err
	}
	v := mockTimeGroupDimension("time_by_hour", "time", 1)
	v.ID = 2
	if err := dictionary.Create([]*models.Dimension{
		{ID:1, Name: "date", FieldName: "date", ValueType: types.ValueTypeString, DataSourceID: 1},
		v,
		{ID:3, Name: "project", FieldName: "project", ValueType: types.ValueTypeString, DataSourceID: 1},
		{ID:4,Name: "sub_project", FieldName: "subproject", ValueType: types.ValueTypeString, DataSourceID: 1},
		{ID:5,Name: "path", FieldName: "path", ValueType: types.ValueTypeString, DataSourceID: 1},
		{ID:6,Name: "project", FieldName: "project", ValueType: types.ValueTypeString, DataSourceID: 2},
		{ID:7,Name: "class_id", FieldName: "class", ValueType: types.ValueTypeInteger, DataSourceID: 2},
		{ID:8,Name: "class_id", FieldName: "id", ValueType: types.ValueTypeInteger, DataSourceID: 3},
		{ID:9,Name: "class_name", FieldName: "name", ValueType: types.ValueTypeString, DataSourceID: 3},
	}); err != nil {
		return err
	}

	if err := dictionary.Create([]*models.DataSet{
		{
			ID:1,
			Name: mockWikiStatDataSet,
			Schema: &models.DataSetSchema{PrimaryID: 1, Secondary: []*models.Secondary{
				{DataSourceID1: 1, DataSourceID2: 2, JoinOn: []*models.JoinOn{{DimensionID1: 3, DimensionID2: 6}}},
				{DataSourceID1: 2, DataSourceID2: 3, JoinOn: []*models.JoinOn{{DimensionID1: 7, DimensionID2: 8}}},
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
	if db, err := client.Get(DataSourceType(), ""); err == nil {
		if err := MockWikiStatData(db); err != nil {
			return err
		}
	}

	return nil
}

// MockQuery1 mock query for normal case
func MockQuery1() *types.Query {
	query := &types.Query{
		DataSetName:  mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:      []string{"hits", "size_sum", "hits_avg", "hits_per_size", "source_avg"},
		Dimensions:   []string{"date", "class_id"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeNotIn, Name: "path", Value: []interface{}{"*"}},
			{OperatorType: types.FilterOperatorTypeIn, Name: "class_id", Value: []interface{}{1, 2, 3, 4}},
		},
		Orders: []*types.OrderBy{
			{Name: "source_sum", Direction: types.OrderDirectionTypeDescending},
		},
		Limit: &types.Limit{Limit: 2, Offset: 1},
	}
	return query
}

func MockQuery1ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 7)
	assert.Equal(t, "date", result.Dimensions[0])
	assert.Equal(t, "source_avg", result.Dimensions[6])
	assert.Len(t, result.Source, 2)
	assert.Len(t, result.Source[0], 7)
	assert.Equal(t, float64(7326), result.Source[0]["size_sum"])
	assert.Equal(t, 0.022113022113022112, result.Source[0]["hits_per_size"])
	assert.Equal(t, 3.700855, result.Source[0]["source_avg"])
}

// MockQuery2 mock query for group by time case
func MockQuery2() *types.Query {
	query := &types.Query{
		DataSetName:  mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:      []string{"hits", "size_sum", "hits_avg", "hits_per_size", "source_avg"},
		Dimensions:   []string{"time_by_hour", "class_id"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeNotIn, Name: "path", Value: []interface{}{"*"}},
			{OperatorType: types.FilterOperatorTypeIn, Name: "class_id", Value: []interface{}{1, 2, 3, 4}},
		},
		Orders: []*types.OrderBy{
			{Name: "time_by_hour", Direction: types.OrderDirectionTypeAscending},
		},
		Limit: &types.Limit{Limit: 10},
	}
	return query
}

func MockQuery2ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 7)
	assert.Equal(t, "time_by_hour", result.Dimensions[0])
	assert.Equal(t, "source_avg", result.Dimensions[6])
	assert.Len(t, result.Source, 8)
	assert.Len(t, result.Source[0], 7)
	assert.Equal(t, "2021-05-06 11:00:00", result.Source[0]["time_by_hour"])
	assert.Equal(t, float64(10086), result.Source[0]["size_sum"])
	assert.Equal(t, 0.013781479278207416, result.Source[0]["hits_per_size"])
	assert.Equal(t, 4.872, result.Source[0]["source_avg"])
}

// MockQuery3 mock query for nested join case
func MockQuery3() *types.Query {
	query := &types.Query{
		DataSetName:  mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:      []string{"source_avg"},
		Dimensions:   []string{"project", "class_name"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeNotIn, Name: "path", Value: []interface{}{"*"}},
			{OperatorType: types.FilterOperatorTypeIn, Name: "project", Value: []interface{}{"city", "school", "music"}},
		},
		Orders: []*types.OrderBy{
			{Name: "project", Direction: types.OrderDirectionTypeDescending},
		},
	}
	return query
}

func MockQuery3ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 3)
	assert.Equal(t, "project", result.Dimensions[0])
	assert.Equal(t, "source_avg", result.Dimensions[2])
	assert.Len(t, result.Source, 3)
	assert.Len(t, result.Source[0], 3)
	assert.Equal(t, "school", result.Source[0]["project"])
	assert.Equal(t, "location", result.Source[0]["class_name"])
	assert.Equal(t, 0.18742, result.Source[0]["source_avg"])
}

// MockQuery4 mock query for nan / inf case
func MockQuery4() *types.Query {
	query := &types.Query{
		DataSetName:  mockWikiStatDataSet,
		TimeInterval: &types.TimeInterval{Name: "date", Start: "2021-05-06", End: "2021-05-08"},
		Metrics:      []string{"hits_per_size"},
		Dimensions:   []string{"class_name"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeIn, Name: "time_by_hour", Value: []interface{}{"2021-05-07 10:00:00"}},
		},
	}
	return query
}

func MockQuery4ResultAssert(t assert.TestingT, result *types.Result) {
	assert.Len(t, result.Dimensions, 2)
	assert.Equal(t, "class_name", result.Dimensions[0])
	assert.Equal(t, "hits_per_size", result.Dimensions[1])
	assert.Len(t, result.Source, 1)
	assert.Len(t, result.Source[0], 2)
	assert.Equal(t, "entertainment", result.Source[0]["class_name"])
	assert.Equal(t, nil, result.Source[0]["hits_per_size"])
}
