package olapsql

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Translator interface {
	Translate(*types.Query) (types.Clause, error)
}

type TranslatorOption struct {
	Adapter IAdapter
	Query   *types.Query
	Set     *models.DataSet
	Root    string
}

func NewTranslator(option *TranslatorOption) (Translator, error) {
	if option.Set == nil {
		option.Set, _ = option.Adapter.GetDataSetByKey(option.Query.DataSetName)
	}
	if option.Root == "" {
		option.Root = option.Set.GetRoot()
	}
	return newBaseTranslator(option)
}

func newBaseTranslator(option *TranslatorOption) (*BaseTranslator, error) {
	tGraph, _ := option.Set.JoinTopologyGraph()

	jBuilder := &JoinTreeBuilder{
		tree:       tGraph.GetTree(option.Root),
		root:       option.Root,
		dictionary: option.Adapter,
	}
	jTree, err := jBuilder.Build()
	if err != nil {
		return nil, err
	}

	mBuilder := &MetricGraphBuilder{
		dbType:     option.Set.DBType,
		dictionary: option.Adapter,
		joinTree:   jTree,
	}
	mGraph, err := mBuilder.Build()
	if err != nil {
		return nil, err
	}

	translator := &BaseTranslator{
		adapter:     option.Adapter,
		set:         option.Set,
		root:        option.Root,
		joinTree:    jTree,
		metricGraph: mGraph,
	}
	return translator, nil
}

type BaseTranslator struct {
	adapter IAdapter
	set     *models.DataSet
	root    string

	joinTree    JoinTree
	metricGraph MetricGraph

	joinedSource []string
}

func (t *BaseTranslator) Translate(query *types.Query) (types.Clause, error) {
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
		DBType:     t.set.DBType,
		Dataset:    t.set.Name,
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

func (t *BaseTranslator) buildMetrics(query *types.Query) ([]*types.Metric, error) {
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		hit, err := t.joinTree.FindMetricByName(v)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(hit.DataSource)
		if err != nil {
			return nil, err
		}
		t.joinedSource = append(t.joinedSource, path...)

		metric, err := t.metricGraph.GetByName(v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (t *BaseTranslator) buildDimensions(query *types.Query) ([]*types.Dimension, error) {
	var dimensions []*types.Dimension
	for _, v := range query.Dimensions {
		dimension, err := t.joinTree.FindDimensionByName(v)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(dimension.DataSource)
		if err != nil {
			return nil, err
		}
		t.joinedSource = append(t.joinedSource, path...)

		d := &types.Dimension{
			Table:     dimension.DataSource,
			Name:      dimension.Name,
			FieldName: dimension.FieldName,
		}

		dimensions = append(dimensions, d)
	}
	return dimensions, nil
}

func (t *BaseTranslator) buildFilters(query *types.Query) ([]*types.Filter, error) {
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

func (t *BaseTranslator) buildOrders(query *types.Query) ([]*types.OrderBy, error) {
	var orders []*types.OrderBy
	for _, v := range query.Orders {
		c, err := t.getColumn(v.Name)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(c.DataSource)
		if err != nil {
			return nil, err
		}
		t.joinedSource = append(t.joinedSource, path...)

		o := &types.OrderBy{
			Name:      c.Statement,
			Direction: v.Direction,
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (t *BaseTranslator) buildJoins() ([]*types.Join, error) {
	var joins []*types.Join
	linq.From(t.joinedSource).Distinct().ToSlice(&t.joinedSource)
	for _, v := range t.joinedSource {
		if v == t.root {
			continue
		}
		join, err := t.getJoin(v)
		if err != nil {
			return nil, err
		}

		var on []*types.JoinOn
		for _, u := range join.JoinOn {
			k1 := fmt.Sprintf("%v.%v", join.DataSource1, u.Dimension1)
			k2 := fmt.Sprintf("%v.%v", join.DataSource2, u.Dimension2)
			d1, _ := t.adapter.GetDimensionByKey(k1)
			d2, _ := t.adapter.GetDimensionByKey(k2)
			on = append(on, &types.JoinOn{Key1: d1.FieldName, Key2: d2.FieldName})
		}

		s1, _ := t.adapter.GetSourceByKey(join.DataSource1)
		s2, _ := t.adapter.GetSourceByKey(join.DataSource2)
		j := &types.Join{
			DataSource1: &types.DataSource{
				Database: s1.Database,
				Name:     s1.Name,
			},
			DataSource2: &types.DataSource{
				Database: s2.Database,
				Name:     s2.Name,
			},
			On: on,
		}

		joins = append(joins, j)
	}
	return joins, nil
}

func (t *BaseTranslator) buildLimit(query *types.Query) (*types.Limit, error) {
	if query.Limit == nil {
		return nil, nil
	}
	return &types.Limit{Limit: query.Limit.Limit, Offset: query.Limit.Offset}, nil
}

func (t *BaseTranslator) buildDataSource() (*types.DataSource, error) {
	source, _ := t.adapter.GetSourceByKey(t.root)
	return &types.DataSource{Database: source.Database, Name: source.Name}, nil
}

func (t *BaseTranslator) getColumn(name string) (*columnStruct, error) {
	metric, err := t.joinTree.FindMetricByName(name)
	if err == nil {
		current, _ := t.metricGraph.GetByName(metric.Name)
		statement, _ := current.Expression()
		return &columnStruct{ValueType: metric.ValueType, Statement: statement, DataSource: metric.DataSource}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	dimension, err := t.joinTree.FindDimensionByName(name)
	if err == nil {
		current := &types.Dimension{
			Table:     dimension.DataSource,
			Name:      dimension.Name,
			FieldName: dimension.FieldName,
		}
		statement, _ := current.Expression()
		return &columnStruct{ValueType: dimension.ValueType, Statement: statement, DataSource: dimension.DataSource}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	return nil, fmt.Errorf("not found filter name %v", name)
}

func (t *BaseTranslator) getJoin(datasource string) (*models.DataSetJoin, error) {
	for _, join := range t.set.Join {
		if join.DataSource2 == datasource {
			return join, nil
		}
	}
	return nil, fmt.Errorf("not found dataset_join data source %v", datasource)
}

func (t *BaseTranslator) treeFilter(in *types.Filter) (*types.Filter, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}

	if !out.OperatorType.IsTree() {
		c, err := t.getColumn(in.Name)
		if err != nil {
			return nil, err
		}
		path, err := t.joinTree.Path(c.DataSource)
		if err != nil {
			return nil, err
		}
		t.joinedSource = append(t.joinedSource, path...)

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
	ValueType  types.ValueType
	Statement  string
	DataSource string
}
