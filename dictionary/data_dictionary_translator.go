package dictionary

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/types"
	"strings"
)

type DictionaryTranslator struct {
	set        *DataSet
	sources    []*DataSource
	metrics    []*Metric
	dimensions []*Dimension

	primaryID      uint64
	joinedSourceID []uint64
	sourceMap      map[uint64]*DataSource
	metricMap      map[uint64]*Metric
	dimensionMap   map[uint64]*Dimension

	joinTree    JoinTree
	metricGraph MetricGraph
}

func (t *DictionaryTranslator) Translate(query *types.Query) (*types.Request, error) {
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
	orders, err := t.buildOrders(query)
	if err != nil {
		return nil, err
	}
	joins, err := t.buildJoins()
	if err != nil {
		return nil, err
	}
	limit, err := t.buildLimit(query)
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
		Orders:     orders,
		Limit:      limit,
		DataSource: datasource,
		Sql:        query.Sql,
	}

	return request, nil
}

func (t *DictionaryTranslator) init() (err error) {
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
	t.joinTree, err = t.buildJoinTree()
	if err != nil {
		return err
	}
	t.metricGraph, err = t.buildMetricGraph()
	if err != nil {
		return err
	}
	return nil
}

func (t *DictionaryTranslator) buildJoinTree() (JoinTree, error) {
	tree, err := t.set.Schema.Tree()
	if err != nil {
		return nil, err
	}

	builder := &JoinTreeBuilder{
		tree: tree,
		root: t.primaryID,
		metrics: t.metrics,
		dimensions: t.dimensions,
		sourceMap: t.sourceMap,
	}
	return builder.Build()
}

func (t *DictionaryTranslator) buildMetricGraph() (MetricGraph, error) {
	builder := &MetricGraphBuilder{sourceMap: t.sourceMap, metricMap: t.metricMap, joinTree: t.joinTree}
	return builder.Build()
}

func (t *DictionaryTranslator) buildMetrics(query *types.Query) ([]*types.Metric, error) {
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		hit, err := t.joinTree.FindMetric(v)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(hit.DataSourceID)
		if err != nil {
			return nil, err
		}
		t.joinedSourceID = append(t.joinedSourceID, path...)

		metric, err := t.metricGraph.GetByName(v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (t *DictionaryTranslator) buildDimensions(query *types.Query) ([]*types.Dimension, error) {
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

func (t *DictionaryTranslator) buildFilters(query *types.Query) ([]*types.Filter, error) {
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

func (t *DictionaryTranslator) buildOrders(query *types.Query) ([]*types.OrderBy, error) {
	var orders []*types.OrderBy
	for _, v := range query.Orders {
		c, err := t.getColumn(v.Name)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(c.DataSourceID)
		if err != nil {
			return nil, err
		}
		t.joinedSourceID = append(t.joinedSourceID, path...)

		o := &types.OrderBy{
			Name:      c.Statement,
			Direction: v.Direction,
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (t *DictionaryTranslator) buildJoins() ([]*types.Join, error) {
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

		j := &types.Join{
			Database1: s1.GetDatabaseName(),
			Database2: s2.GetDatabaseName(),
			Table1:    s1.GetTableName(),
			Table2:    s2.GetTableName(),
			On:        on,
		}

		joins = append(joins, j)
	}
	return joins, nil
}

func (t *DictionaryTranslator) buildLimit(query *types.Query) (*types.Limit, error) {
	if query.Limit == nil {
		return nil, nil
	}
	return &types.Limit{Limit: query.Limit.Limit, Offset: query.Limit.Offset}, nil
}

func (t *DictionaryTranslator) buildDataSource() (*types.DataSource, error) {
	source := t.sourceMap[t.primaryID]
	return &types.DataSource{Type: source.Type, Name: source.Name}, nil
}

func (t *DictionaryTranslator) getColumn(name string) (*columnStruct, error) {
	metric, err := t.joinTree.FindMetric(name)
	if err == nil {
		current, _ := t.metricGraph.GetByID(metric.ID)
		statement, _ := current.Expression()
		return &columnStruct{ValueType: metric.ValueType, Statement: statement, DataSourceID: metric.DataSourceID}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	dimension, err := t.joinTree.FindDimension(name)
	if err == nil {
		source := t.sourceMap[dimension.DataSourceID]
		current := &types.Dimension{
			Table:     source.GetTableName(),
			Name:      dimension.Name,
			FieldName: dimension.FieldName,
		}
		statement, _ := current.Expression()
		return &columnStruct{ValueType: dimension.ValueType, Statement: statement, DataSourceID: dimension.DataSourceID}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	return nil, fmt.Errorf("not found filter name %v", name)
}

func (t *DictionaryTranslator) getDataSource(id uint64) (*models.DataSource, error) {
	if v, ok := t.sourceMap[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found data source id %v", id)
}

func (t *DictionaryTranslator) getSecondary(id uint64) (*models.Secondary, error) {
	if v, ok := t.secondaryMap[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found secondary data source id %v", id)
}

func (t *DictionaryTranslator) treeFilter(in *types.Filter) (*types.Filter, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}

	if !out.OperatorType.IsTree() {
		c, err := t.getColumn(in.Name)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(c.DataSourceID)
		if err != nil {
			return nil, err
		}
		t.joinedSourceID = append(t.joinedSourceID, path...)

		out.ValueType = c.ValueType
		out.Name = c.Statement
		return out, nil
	}

	for _, v := range in.Children {
		child, err := t.treeFilter(v)
		if err != nil {
			return nil, err
		}
		out.Children = append(out.Children, child)
	}

	return out, nil
}

type columnStruct struct {
	ValueType    types.ValueType `json:"value_type"`
	Statement    string          `json:"statement"`
	DataSourceID uint64          `json:"data_source_id"`
}

