package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type Translator interface {
	Translate(interface{}) (interface{}, error)
}

type translator struct {
	db         *gorm.DB
	set        *models.DataSet
	sources    []*models.DataSource
	metrics    []*models.Metric
	dimensions []*models.Dimension

	primaryID uint64
	joinedID  []uint64
	sourceMap map[uint64]*models.DataSource
	metricMap map[string][]*models.Metric
	dimensionMap map[string][]*models.Dimension
	secondaryMap map[uint64]*models.Secondary
}

func (t *translator) Translate(in interface{}) (interface{}, error) {
	return nil, nil
}
// func (t *translator) Translate(in interface{}) (interface{}, error) {
// 	query, ok := in.(*types.Query)
// 	if !ok {
// 		return nil, fmt.Errorf("not supported type %T", in)
// 	}
//
// 	if err := t.init(); err != nil {
// 		return nil, err
// 	}
//
// 	var joined []uint64
// 	var metrics []*types.Metric
// 	for _, v := range query.Metrics {
// 		m, err := t.getMetric(v)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if m.DataSourceID != primaryID {
// 			joined = append(joined, m.DataSourceID)
// 		}
//
// 		metrics = append(metrics, &types.Metric{
// 			Type: m.Type,
// 			FieldName: m.FieldName,
// 			Name: m.Name,
// 		})
// 	}
//
// 	for _, v := range query.Dimensions {
// 		m, err := t.getMetric(v)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if m.DataSourceID != primaryID {
// 			joined = append(joined, m.DataSourceID)
// 		}
//
// 		metrics = append(metrics, &types.Metric{
// 			Type: m.Type,
// 			FieldName: m.FieldName,
// 			Name: m.Name,
// 		})
// 	}
//
// 	for _, v := range query.Filters {
// 		v.
// 		m, err := t.getMetric(v)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if m.DataSourceID != primaryID {
// 			joined = append(joined, m.DataSourceID)
// 		}
//
// 		metrics = append(metrics, &types.Metric{
// 			Type: m.Type,
// 			FieldName: m.FieldName,
// 			Name: m.Name,
// 		})
// 	}
//
// 	request := &types.Request{
// 		DataSource: &types.DataSource{Type: source.Type, Name: source.Name},
// 		Metrics: metrics,
// 	}
//
// 	return request, nil
// }

func (t *translator) init() error {
	t.primaryID = t.set.Schema.PrimaryID
	t.sourceMap = map[uint64]*models.DataSource{}
	for _, v := range t.sources {
		t.sourceMap[v.ID] = v
	}
	t.metricMap = map[string][]*models.Metric{}
	for _, v := range t.metrics {
		t.metricMap[v.Name] = append(t.metricMap[v.Name], v)
	}
	t.dimensionMap = map[string][]*models.Dimension{}
	for _, v := range t.dimensions {
		t.dimensionMap[v.Name] = append(t.dimensionMap[v.Name], v)
	}
	t.secondaryMap = map[uint64]*models.Secondary{}
	for _, v := range t.set.Schema.Secondary {
		t.secondaryMap[v.DataSourceID] = v
	}
	return nil
}

func (t *translator) buildMetrics(query *types.Query) ([]*types.Metric, error) {
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		m, err := t.getMetric(v)
		if err != nil {
			return nil, err
		}

		if m.DataSourceID != t.primaryID {
			t.joinedID = append(t.joinedID, m.DataSourceID)
		}

		metrics = append(metrics, &types.Metric{
			Type: m.Type,
			FieldName: m.FieldName,
			Name: m.Name,
		})
	}
	return metrics, nil
}

func (t *translator) buildDimensions(query *types.Query) ([]*types.Dimension, error) {
	return nil, nil
}

func (t *translator) buildFilters(query *types.Query) ([]*types.Filter, error) {
	return nil, nil
}

func (t *translator) buildJoins(joined []uint64) ([]*types.Join, error) {
	return nil, nil
}

func (t *translator) getDataSource(id uint64) (*models.DataSource, error) {
	if v, ok := t.sourceMap[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found data source id %v", id)
}

func (t *translator) getMetric(name string) (*models.Metric, error) {
	if v, ok := t.metricMap[name]; ok {
		if len(v) >= 2 {
			return nil, fmt.Errorf("duplicate metric name %v", name)
		}
		return v[0], nil
	}

	return nil, fmt.Errorf("not found metric name %v", name)
}

func (t *translator) getDimension(name string) (*models.Dimension, error) {
	if v, ok := t.dimensionMap[name]; ok {
		if len(v) >= 2 {
			return nil, fmt.Errorf("duplicate dimension name %v", name)
		}
		return v[0], nil
	}

	return nil, fmt.Errorf("not found dimension name %v", name)
}
