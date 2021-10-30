package olapsql

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Translator interface {
	GetAdapter() IAdapter
	GetJoinTree() JoinTree
	GetMetricGraph() MetricGraph
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

	tGraph, _ := GetDependencyTree(adapter, option.Current)

	jBuilder := &JoinTreeBuilder{
		tree:       tGraph.GetTree(option.Current),
		root:       option.Current,
		dictionary: adapter,
	}
	jTree, err := jBuilder.Build()
	if err != nil {
		return nil, err
	}

	mBuilder := &MetricGraphBuilder{
		dbType:     option.DBType,
		dictionary: adapter,
		joinTree:   jTree,
	}
	mGraph, err := mBuilder.Build()
	if err != nil {
		return nil, err
	}

	translator := &normalTranslator{
		adapter:     adapter,
		query:       option.Query,
		dBType:      option.DBType,
		current:     option.Current,
		joinTree:    jTree,
		metricGraph: mGraph,
	}
	return translator, nil
}

func newDirectSqlTranslator(option *TranslatorOption) (*directSqlTranslator, error) {
	translator := &directSqlTranslator{dBType: option.DBType}
	return translator, nil
}

type columnStruct struct {
	ValueType  types.ValueType
	Statement  string
	DataSource string
}

func getColumn(translator Translator, name string) (*columnStruct, error) {
	joinTree := translator.GetJoinTree()
	metricGraph := translator.GetMetricGraph()
	metric, err := joinTree.FindMetricByName(name)
	if err == nil {
		current, _ := metricGraph.GetByName(metric.Name)
		statement, _ := current.Expression()
		return &columnStruct{ValueType: metric.ValueType, Statement: statement, DataSource: metric.DataSource}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate filter name %v", name)
	}

	dimension, err := joinTree.FindDimensionByName(name)
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

func buildMetrics(translator Translator, query *types.Query) ([]*types.Metric, error) {
	joinTree := translator.GetJoinTree()
	metricGraph := translator.GetMetricGraph()
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		_, err := joinTree.FindMetricByName(v)
		if err != nil {
			return nil, err
		}

		metric, err := metricGraph.GetByName(v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func buildDimensions(translator Translator, query *types.Query) ([]*types.Dimension, error) {
	joinTree := translator.GetJoinTree()
	var dimensions []*types.Dimension
	for _, v := range query.Dimensions {
		hit, err := joinTree.FindDimensionByName(v)
		if err != nil {
			return nil, err
		}

		d := &types.Dimension{
			Table:     hit.DataSource,
			Name:      hit.Name,
			FieldName: hit.FieldName,
		}
		dimensions = append(dimensions, d)
	}
	return dimensions, nil
}

func buildOneFilter(translator Translator, in *types.Filter) (*types.Filter, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}

	if !out.OperatorType.IsTree() {
		c, err := getColumn(translator, in.Name)
		if err != nil {
			return nil, err
		}
		out.ValueType = c.ValueType
		out.Table = c.DataSource
		out.Name = c.Statement
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
		c, err := getColumn(translator, v.Name)
		if err != nil {
			return nil, err
		}

		o := &types.OrderBy{
			Table:     c.DataSource,
			Name:      c.Statement,
			Direction: v.Direction,
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
	adapter IAdapter
	query   *types.Query
	dBType  types.DBType
	current string

	joinTree    JoinTree
	metricGraph MetricGraph
}

func (n *normalTranslator) GetAdapter() IAdapter {
	return n.adapter
}

func (n *normalTranslator) GetJoinTree() JoinTree {
	return n.joinTree
}

func (n *normalTranslator) GetMetricGraph() MetricGraph {
	return n.metricGraph
}

func (n *normalTranslator) Translate(query *types.Query) (types.Clause, error) {
	clause, err := buildNormalClause(n, query)
	if err != nil {
		return nil, err
	}
	clause.DataSource, clause.Joins, err = n.buildDataSourcesAndJoins()
	if err != nil {
		return nil, err
	}
	clause.DBType = n.dBType
	clause.Dataset = query.DataSetName
	return clause, nil
}

func (n *normalTranslator) buildDataSourcesAndJoins() (sources []*types.DataSource, joins []*types.Join, err error) {
	toFn := func(in *models.DataSource) *types.DataSource {
		return &types.DataSource{Database: in.Database, Name: in.Name, AliasName: in.Alias, Type: in.Type}
	}

	source, _ := n.adapter.GetSourceByKey(n.current)
	switch source.Type {
	case types.DataSourceTypeFact:
		sources = append(sources, toFn(source))
		return
	case types.DataSourceTypeFactDimensionJoin:
	case types.DataSourceTypeMergedJoin:
		return
	default:
		err = fmt.Errorf("can't use datasource type=%v as dateset's datasource", source.Type)
		return
	}

	for _, v := range t.joinedSource {
		if v == t.current {
			continue
		}
		join, err := t.getJoin(v)
		if err != nil {
			return nil, err
		}

		ds1, dl1, ds2, dl2 := join.Get1().DataSource, join.Get1().Dimension, join.Get2().DataSource, join.Get2().Dimension
		var on []*types.JoinOn
		for i := 0; i <= len(dl1); i++ {
			k1 := fmt.Sprintf("%v.%v", ds1, dl1)
			k2 := fmt.Sprintf("%v.%v", ds2, dl2)
			d1, _ := t.adapter.GetDimensionByKey(k1)
			d2, _ := t.adapter.GetDimensionByKey(k2)
			on = append(on, &types.JoinOn{Key1: d1.FieldName, Key2: d2.FieldName})
		}

		s1, _ := t.adapter.GetSourceByKey(ds1)
		s2, _ := t.adapter.GetSourceByKey(ds2)
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

type directSqlTranslator struct {
	dBType types.DBType
}

func (d *directSqlTranslator) GetAdapter() IAdapter {
	panic("implement me")
}

func (d *directSqlTranslator) GetJoinTree() JoinTree {
	panic("implement me")
}

func (d *directSqlTranslator) GetMetricGraph() MetricGraph {
	panic("implement me")
}

func (d *directSqlTranslator) Translate(query *types.Query) (types.Clause, error) {
	clause := &types.SqlClause{Sql: query.Sql}
	clause.DBType = d.dBType
	clause.Dataset = query.DataSetName
	return clause, nil
}
