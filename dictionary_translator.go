package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
)

type Translator interface {
	GetCurrent() string
	GetAdapter() IAdapter
	GetDependencyGraph() DependencyGraph
	Translate(*types.Query) (types.Clause, error)
}

type TranslatorOption struct {
	Adapter IAdapter
	Query   *types.Query
	DBType  types.DBType
	Current string
}

func NewTranslator(option *TranslatorOption) (Translator, error) {
	if option.Query.Sql != "" {
		return newDirectSqlTranslator(option)
	}
	return newNormalTranslator(option)
}

func newNormalTranslator(option *TranslatorOption) (*normalTranslator, error) {
	adapter, err := option.Adapter.BuildDataSourceAdapter(option.Current)
	if err != nil {
		return nil, err
	}

	dBuilder := &DependencyGraphBuilder{
		dbType:     option.DBType,
		dictionary: adapter,
		current:    option.Current,
	}
	dGraph, err := dBuilder.Build()
	if err != nil {
		return nil, err
	}

	translator := &normalTranslator{
		option:  option,
		adapter: adapter,
		dBType:  option.DBType,
		current: option.Current,
		dGraph:  dGraph,
	}
	return translator, nil
}

func newDirectSqlTranslator(option *TranslatorOption) (*directSqlTranslator, error) {
	translator := &directSqlTranslator{dBType: option.DBType}
	return translator, nil
}

func getColumn(translator Translator, table, name string) (*columnStruct, error) {
	key := fmt.Sprintf("%v.%v", table, name)
	dGraph := translator.GetDependencyGraph()
	if metric, err := dGraph.GetMetric(key); err == nil {
		return &columnStruct{FieldProperty: types.FieldPropertyMetric, Metric: metric}, nil
	}
	if dimension, err := dGraph.GetDimension(key); err == nil {
		return &columnStruct{FieldProperty: types.FieldPropertyDimension, Dimension: dimension}, nil
	}
	return nil, fmt.Errorf("not found column name %v", name)
}

func buildMetrics(translator Translator, query *types.Query) ([]*types.Metric, error) {
	dGraph := translator.GetDependencyGraph()
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		key := fmt.Sprintf("%v.%v", translator.GetCurrent(), v)
		metric, err := dGraph.GetMetric(key)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func buildDimensions(translator Translator, query *types.Query) ([]*types.Dimension, error) {
	dGraph := translator.GetDependencyGraph()
	var dimensions []*types.Dimension
	for _, v := range query.Dimensions {
		key := fmt.Sprintf("%v.%v", translator.GetCurrent(), v)
		dimension, err := dGraph.GetDimension(key)
		if err != nil {
			return nil, err
		}
		dimensions = append(dimensions, dimension)
	}
	return dimensions, nil
}

func buildOneFilter(translator Translator, in *types.Filter) (*types.Filter, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}

	if !out.OperatorType.IsTree() {
		c, err := getColumn(translator, translator.GetCurrent(), in.Name)
		if err != nil {
			return nil, err
		}
		out.FieldProperty = c.FieldProperty
		out.ValueType = c.GetValueType()
		out.Table = c.GetTable()
		out.Name = c.GetExpression()
		return out, nil
	}

	for _, v := range in.Children {
		child, err := buildOneFilter(translator, v)
		if err != nil {
			return nil, err
		}
		out.Children = append(out.Children, child)
	}

	return out, nil
}

func buildFilters(translator Translator, query *types.Query) ([]*types.Filter, error) {
	var filters []*types.Filter
	for _, v := range query.Filters {
		filter, err := buildOneFilter(translator, v)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func buildOrders(translator Translator, query *types.Query) ([]*types.OrderBy, error) {
	var orders []*types.OrderBy
	for _, v := range query.Orders {
		c, err := getColumn(translator, translator.GetCurrent(), v.Name)
		if err != nil {
			return nil, err
		}
		o := &types.OrderBy{
			Table:         c.GetTable(),
			Name:          c.GetExpression(),
			FieldProperty: c.FieldProperty,
			Direction:     v.Direction,
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func buildLimit(_ Translator, query *types.Query) (*types.Limit, error) {
	if query.Limit == nil {
		return nil, nil
	}
	return &types.Limit{Limit: query.Limit.Limit, Offset: query.Limit.Offset}, nil
}

func buildNormalClause(translator Translator, query *types.Query) (*types.NormalClause, error) {
	var err error
	clause := &types.NormalClause{}
	clause.Metrics, err = buildMetrics(translator, query)
	if err != nil {
		return nil, err
	}
	clause.Dimensions, err = buildDimensions(translator, query)
	if err != nil {
		return nil, err
	}
	clause.Filters, err = buildFilters(translator, query)
	if err != nil {
		return nil, err
	}
	clause.Orders, err = buildOrders(translator, query)
	if err != nil {
		return nil, err
	}
	clause.Limit, err = buildLimit(translator, query)
	if err != nil {
		return nil, err
	}
	return clause, nil
}

type normalTranslator struct {
	option  *TranslatorOption
	adapter IAdapter
	dBType  types.DBType
	current string

	dGraph DependencyGraph
}

func (n *normalTranslator) GetCurrent() string {
	return n.current
}

func (n *normalTranslator) GetAdapter() IAdapter {
	return n.adapter
}

func (n *normalTranslator) GetDependencyGraph() DependencyGraph {
	return n.dGraph
}

func (n *normalTranslator) Translate(query *types.Query) (types.Clause, error) {
	clause, err := buildNormalClause(n, query)
	if err != nil {
		return nil, err
	}
	splitter, err := NewNormalClauseSplitter(n, clause, query, n.dBType)
	if err != nil {
		return nil, err
	}
	if err = splitter.Run(); err != nil {
		return nil, err
	}
	clause.DBType = n.dBType
	clause.Dataset = query.DataSetName
	return clause, nil
}

type directSqlTranslator struct {
	dBType types.DBType
}

func (d *directSqlTranslator) GetCurrent() string {
	panic("implement me")
}

func (d *directSqlTranslator) GetAdapter() IAdapter {
	panic("implement me")
}

func (d *directSqlTranslator) GetDependencyGraph() DependencyGraph {
	panic("implement me")
}

func (d *directSqlTranslator) Translate(query *types.Query) (types.Clause, error) {
	clause := &types.SqlClause{Sql: query.Sql}
	clause.DBType = d.dBType
	clause.Dataset = query.DataSetName
	return clause, nil
}
