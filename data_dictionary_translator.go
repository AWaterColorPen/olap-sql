package olapsql

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type dataDictionaryTranslator struct {
	db         *gorm.DB
	set        *models.DataSet
	sources    []*models.DataSource
	metrics    []*models.Metric
	dimensions []*models.Dimension

	primaryID        uint64
	joinedSourceID   []uint64
	sourceMap        map[uint64]*models.DataSource
	metricMap        map[uint64]*models.Metric
	dimensionMap     map[uint64]*models.Dimension
	secondaryMap     map[uint64]*models.Secondary
	metricNameMap    map[string][]*models.Metric
	dimensionNameMap map[string][]*models.Dimension

	joinTree JoinTree
}

func (t *dataDictionaryTranslator) Translate(query *types.Query) (*types.Request, error) {
	if err := t.init(); err != nil {
		return nil, err
	}

	metrics, err := t.buildMetrics(query)
	if err != nil {
		return nil, err
	}
	dimensions, err := t.buildDimensions(query)
	if err != nil {
		return nil, err
	}
	filters, err := t.buildFilters(query)
	if err != nil {
		return nil, err
	}
	joins, err := t.buildJoins()
	if err != nil {
		return nil, err
	}
	datasource, err := t.buildDataSource()
	if err != nil {
		return nil, err
	}
	request := &types.Request{
		Metrics:    metrics,
		Dimensions: dimensions,
		Filters:    filters,
		Joins:      joins,
		DataSource: datasource,
	}

	return request, nil
}

func (t *dataDictionaryTranslator) init() error {
	t.primaryID = t.set.Schema.PrimaryID
	t.sourceMap = map[uint64]*models.DataSource{}
	for _, v := range t.sources {
		t.sourceMap[v.ID] = v
	}
	t.metricMap = Metrics(t.metrics).IdIndex()
	t.dimensionMap = Dimensions(t.dimensions).IdIndex()
	t.secondaryMap = map[uint64]*models.Secondary{}
	for _, v := range t.set.Schema.Secondary {
		t.secondaryMap[v.DataSourceID2] = v
	}
	t.metricNameMap = map[string][]*models.Metric{}
	for _, v := range t.metrics {
		t.metricNameMap[v.Name] = append(t.metricNameMap[v.Name], v)
	}
	t.dimensionNameMap = map[string][]*models.Dimension{}
	for _, v := range t.dimensions {
		t.dimensionNameMap[v.Name] = append(t.dimensionNameMap[v.Name], v)
	}
	tree, err := t.buildJoinTree()
	if err != nil {
		return err
	}
	t.joinTree = tree
	return nil
}

func (t *dataDictionaryTranslator) buildJoinTree() (JoinTree, error) {
	err := t.set.Schema.Valid(t.db)
	if err != nil {
		return nil, err
	}

	tree, err := t.set.Schema.Tree()
	if err != nil {
		return nil, err
	}

	builder := &JoinTreeBuilder{tree: tree, root: t.primaryID, sourceMap: t.sourceMap}
	return builder.Build()
}

func (t *dataDictionaryTranslator) buildMetrics(query *types.Query) ([]*types.Metric, error) {
	originMetrics := map[string]bool{}
	for _, v := range query.Metrics {
		originMetrics[v] = true
	}

	queue, err := t.sortMetrics(query)
	if err != nil {
		return nil, err
	}

	var metrics []*types.Metric
	metricsMap := map[uint64]*types.Metric{}
	for _, v := range queue {
		m := t.metricMap[v]
		source := t.sourceMap[m.DataSourceID]

		tm := &types.Metric{
			Type:      m.Type,
			Table:     source.GetTableName(),
			Name:      m.Name,
			FieldName: m.FieldName,
		}

		if m.Composition != nil {
			for _, u := range m.Composition.MetricID {
				mm := metricsMap[u]
				tm.Children = append(tm.Children, mm)
			}
		}

		metricsMap[v] = tm
		if _, ok := originMetrics[m.Name]; ok {
			metrics = append(metrics, tm)
		}
	}
	return metrics, nil
}

func (t *dataDictionaryTranslator) buildDimensions(query *types.Query) ([]*types.Dimension, error) {
	var dimensions []*types.Dimension
	for _, v := range query.Dimensions {
		dimension, err := t.joinTree.FindDimension(v)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(dimension.DataSourceID)
		if err != nil {
			return nil, err
		}
		t.joinedSourceID = append(t.joinedSourceID, path...)

		source := t.sourceMap[dimension.DataSourceID]
		d := &types.Dimension{
			Table:     source.GetTableName(),
			Name:      dimension.Name,
			FieldName: dimension.FieldName,
		}

		dimensions = append(dimensions, d)
	}
	return dimensions, nil
}

func (t *dataDictionaryTranslator) buildFilters(query *types.Query) ([]*types.Filter, error) {
	var filters []*types.Filter
	for _, v := range query.Filters {
		filter, err := t.treeFilter(v)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func (t *dataDictionaryTranslator) buildJoins() ([]*types.Join, error) {
	var joins []*types.Join
	linq.From(t.joinedSourceID).Distinct().ToSlice(&t.joinedSourceID)
	for _, v := range t.joinedSourceID {
		if v == t.primaryID {
			continue
		}
		s, err := t.getSecondary(v)
		if err != nil {
			return nil, err
		}

		var on []*types.JoinOn
		for _, u := range s.JoinOn {
			d1, ok1 := t.dimensionMap[u.DimensionID1]
			if !ok1 {
				return nil, fmt.Errorf("not found dimension id %v", u.DimensionID1)
			}
			d2, ok2 := t.dimensionMap[u.DimensionID2]
			if !ok2 {
				return nil, fmt.Errorf("not found dimension id %v", u.DimensionID2)
			}
			on = append(on, &types.JoinOn{Key1: d1.FieldName, Key2: d2.FieldName})
		}

		s1 := t.sourceMap[s.DataSourceID1]
		s2 := t.sourceMap[s.DataSourceID2]

		tj := &types.Join{
			Database1: s1.GetDatabaseName(),
			Database2: s2.GetDatabaseName(),
			Table1:    s1.GetTableName(),
			Table2:    s2.GetTableName(),
			On:        on,
		}

		joins = append(joins, tj)
	}
	return joins, nil
}

func (t *dataDictionaryTranslator) buildDataSource() (*types.DataSource, error) {
	source := t.sourceMap[t.primaryID]
	return &types.DataSource{Type: source.Type, Name: source.Name}, nil
}

func (t *dataDictionaryTranslator) getFilter(name string) (*filterStruct, error) {
	m, err := t.joinTree.FindMetric(name)
	if err == nil {
		return &filterStruct{ValueType: m.ValueType, Name: m.FieldName, DataSourceID: m.DataSourceID}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	d, err := t.joinTree.FindDimension(name)
	if err == nil {
		return &filterStruct{ValueType: d.ValueType, Name: d.FieldName, DataSourceID: d.DataSourceID}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	return nil, fmt.Errorf("not found filter name %v", name)
}

func (t *dataDictionaryTranslator) getDataSource(id uint64) (*models.DataSource, error) {
	if v, ok := t.sourceMap[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found data source id %v", id)
}

func (t *dataDictionaryTranslator) getSecondary(id uint64) (*models.Secondary, error) {
	if v, ok := t.secondaryMap[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found secondary data source id %v", id)
}

func (t *dataDictionaryTranslator) sortMetrics(query *types.Query) ([]uint64, error) {
	inDegree := map[uint64]int{}
	graph := map[uint64][]uint64{}

	for _, v := range query.Metrics {
		metric, err := t.joinTree.FindMetric(v)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(metric.DataSourceID)
		if err != nil {
			return nil, err
		}
		t.joinedSourceID = append(t.joinedSourceID, path...)

		inDegree[metric.ID] = 0
		if metric.Composition != nil {
			for _, u := range metric.Composition.MetricID {
				graph[u] = append(graph[u], metric.ID)
				inDegree[metric.ID]++
				if _, ok := inDegree[u]; !ok {
					inDegree[u] = 0
				}
			}
		}
	}

	var queue []uint64
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}

	for i := 0; i < len(queue); i++ {
		node := queue[i]
		for _, v := range graph[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	return queue, nil
}

func (t *dataDictionaryTranslator) treeFilter(in *types.Filter) (*types.Filter, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}

	if !out.OperatorType.IsTree() {
		f, err := t.getFilter(in.Name)
		if err != nil {
			return nil, err
		}
		if f.DataSourceID != t.primaryID {
			t.joinedSourceID = append(t.joinedSourceID, f.DataSourceID)
		}

		source := t.sourceMap[f.DataSourceID]
		out.ValueType = f.ValueType
		out.Name = f.Name
		out.Table = source.Name
		return out, nil
	}

	for _, v := range in.Children {
		current, err := t.treeFilter(v)
		if err != nil {
			return nil, err
		}
		out.Children = append(out.Children, current)
	}

	return out, nil
}

type filterStruct struct {
	ValueType    types.ValueType `json:"value_type"`
	Name         string          `json:"name"`
	DataSourceID uint64          `json:"data_source_id"`
}
