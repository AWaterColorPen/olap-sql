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
		option:      option,
		adapter:     adapter,
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
	ValueType     types.ValueType
	FieldProperty types.FieldProperty
	Statement     string
	DataSource    string
}

func getColumn(translator Translator, name string) (*columnStruct, error) {
	jTree := translator.GetJoinTree()
	mGraph := translator.GetMetricGraph()
	metric, err := jTree.FindMetricByName(name)
	if err == nil {
		current, _ := mGraph.GetByName(metric.Name)
		statement, _ := current.Expression()
		return &columnStruct{ValueType: metric.ValueType, FieldProperty: types.FieldPropertyMetric, Statement: statement, DataSource: metric.DataSource}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate column name %v", name)
	}

	dimension, err := jTree.FindDimensionByName(name)
	if err == nil {
		current := &types.Dimension{
			Table:       dimension.DataSource,
			Name:        dimension.Name,
			FieldName:   dimension.FieldName,
			Type:        dimension.Type,
			Composition: dimension.Composition,
		}
		statement, _ := current.Expression()
		return &columnStruct{ValueType: dimension.ValueType, FieldProperty: types.FieldPropertyDimension, Statement: statement, DataSource: dimension.DataSource}, nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return nil, fmt.Errorf("duplicate column name %v", name)
	}

	return nil, fmt.Errorf("not found column name %v", name)
}

func buildMetrics(translator Translator, query *types.Query) ([]*types.Metric, error) {
	jTree := translator.GetJoinTree()
	mGraph := translator.GetMetricGraph()
	var metrics []*types.Metric
	for _, v := range query.Metrics {
		_, err := jTree.FindMetricByName(v)
		if err != nil {
			return nil, err
		}

		metric, err := mGraph.GetByName(v)
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
			Table:       hit.DataSource,
			Name:        hit.Name,
			FieldName:   hit.FieldName,
			Type:        hit.Type,
			Composition: hit.Composition,
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
		out.FieldProperty = c.FieldProperty
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
			Table:         c.DataSource,
			Name:          c.Statement,
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

func getHitDatasourceFromMetric(metric *types.Metric) []string {
	hit := []string{metric.Table}
	for _, child := range metric.Children {
		hit = append(hit, getHitDatasourceFromMetric(child)...)
	}
	hit = append(hit, getHitDatasourceFromFilter(metric.Filter)...)
	return hit
}

func getHitDatasourceFromDimension(dimension *types.Dimension) []string {
	return []string{dimension.Table}
}

func getHitDatasourceFromFilter(filter *types.Filter) []string {
	if filter == nil {
		return nil
	}
	var hit []string
	if len(filter.Table) > 0 {
		hit = append(hit, filter.Table)
	}
	for _, child := range filter.Children {
		hit = append(hit, getHitDatasourceFromFilter(child)...)
	}
	return hit
}

func getHitDatasourceFromOrderBy(order *types.OrderBy) []string {
	return []string{order.Table}
}

func getHitDatasourceKey(clause *types.NormalClause) []string {
	var hit []string
	for _, m := range clause.Metrics {
		hit = append(hit, getHitDatasourceFromMetric(m)...)
	}
	for _, d := range clause.Dimensions {
		hit = append(hit, getHitDatasourceFromDimension(d)...)
	}
	for _, f := range clause.Filters {
		hit = append(hit, getHitDatasourceFromFilter(f)...)
	}
	for _, o := range clause.Orders {
		hit = append(hit, getHitDatasourceFromOrderBy(o)...)
	}
	linq.From(hit).Distinct().ToSlice(&hit)
	return hit
}

func getHitDatasource(translator Translator, clause *types.NormalClause) ([]*types.DataSource, error) {
	adapter := translator.GetAdapter()
	hitKey := getHitDatasourceKey(clause)
	var source []*types.DataSource
	for _, hit := range hitKey {
		s, err := adapter.GetSourceByKey(hit)
		if err != nil {
			return nil, err
		}
		ss := &types.DataSource{
			Database:  s.Database,
			Name:      s.Name,
			AliasName: s.Alias,
			Type:      s.Type,
		}
		source = append(source, ss)
	}
	return nil, nil
}

// func MD(clause *types.NormalClause, candidate map[string]*types.DataSource) (map[*types.DataSource]*types.Query, error) {
// 	for
//
// }

func buildDimensionJoin(translator Translator, source *models.DataSource, hitMap map[string]*types.DataSource) []*types.Join {
	adapter := translator.GetAdapter()
	var joins []*types.Join
	for _, v := range source.DimensionJoin {
		s1, ok1 := hitMap[v.Get1().DataSource]
		s2, ok2 := hitMap[v.Get2().DataSource]
		if !ok1 || !ok2 {
			continue
		}
		ds1, dl1, ds2, dl2 := v.Get1().DataSource, v.Get1().Dimension, v.Get2().DataSource, v.Get2().Dimension
		var on []*types.JoinOn
		for i := 0; i <= len(dl1); i++ {
			k1 := fmt.Sprintf("%v.%v", ds1, dl1[i])
			k2 := fmt.Sprintf("%v.%v", ds2, dl2[i])
			d1, _ := adapter.GetDimensionByKey(k1)
			d2, _ := adapter.GetDimensionByKey(k2)
			on = append(on, &types.JoinOn{Key1: d1.FieldName, Key2: d2.FieldName})
		}

		j := &types.Join{DataSource1: s1, DataSource2: s2, On: on}
		joins = append(joins, j)
	}
	return joins
}

func buildMergedJoin(translator Translator, source *models.DataSource, hitMap map[string]*types.DataSource) []*types.Join {
	adapter := translator.GetAdapter()
	var joins []*types.Join
	for i := 2; i < len(source.MergedJoin); i++ {
		s1, ok1 := hitMap[source.MergedJoin[1].DataSource]
		s2, ok2 := hitMap[source.MergedJoin[i].DataSource]
		if !ok1 || !ok2 {
			continue
		}
		ds1, dl1 := source.MergedJoin[1].DataSource, source.MergedJoin[1].Dimension
		ds2, dl2 := source.MergedJoin[i].DataSource, source.MergedJoin[i].Dimension
		var on []*types.JoinOn
		for j := 0; j <= len(dl1); i++ {
			k1 := fmt.Sprintf("%v.%v", ds1, dl1[i])
			k2 := fmt.Sprintf("%v.%v", ds2, dl2[i])
			d1, _ := adapter.GetDimensionByKey(k1)
			d2, _ := adapter.GetDimensionByKey(k2)
			on = append(on, &types.JoinOn{Key1: d1.FieldName, Key2: d2.FieldName})
		}
		j := &types.Join{DataSource1: s1, DataSource2: s2, On: on}
		joins = append(joins, j)
	}
	return joins
}

func buildJoins(translator Translator, source *models.DataSource, hit []*types.DataSource) ([]*types.Join, error) {
	hitMap := map[string]*types.DataSource{}
	for _, v := range hit {
		hitMap[v.Name] = v
	}
	var joins []*types.Join
	joins = append(joins, buildDimensionJoin(translator, source, hitMap)...)
	joins = append(joins, buildMergedJoin(translator, source, hitMap)...)
	return joins, nil
}

type normalTranslator struct {
	option  *TranslatorOption
	adapter IAdapter
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
	clause.DataSource, clause.Joins, err = n.buildDataSourcesAndJoins(clause)
	if err != nil {
		return nil, err
	}
	clause.DBType = n.dBType
	clause.Dataset = query.DataSetName
	return clause, nil
}

func (n *normalTranslator) buildDataSourcesAndJoins(clause *types.NormalClause) (sources []*types.DataSource, joins []*types.Join, err error) {
	toFn := func(in *models.DataSource) *types.DataSource {
		return &types.DataSource{Database: in.Database, Name: in.Name, AliasName: in.Alias, Type: in.Type}
	}

	source, _ := n.adapter.GetSourceByKey(n.current)

	switch source.Type {
	case types.DataSourceTypeFact:
		sources = append(sources, toFn(source))
		return
	case types.DataSourceTypeFactDimensionJoin:
		sources, err = getHitDatasource(n, clause)
		if err != nil {
			return
		}
		joins, err = buildJoins(n, source, sources)
		return
	case types.DataSourceTypeMergedJoin:
		sources, err = getHitDatasource(n, clause)
		if err != nil {
			return
		}
		joins, err = buildJoins(n, source, sources)
		if err != nil {
			return
		}
		var splitter *normalClauseSplitter
		splitter, err = NewNormalClauseSplitter(sources)
		if err != nil {
			return
		}
		mq, ee := splitter.Split(clause)
		if ee != nil {
			return
		}
		for k, v := range mq {
			o := &TranslatorOption{
				Adapter: n.option.Adapter,
				Query:   v,
				DBType:  n.dBType,
				Current: n.current,
			}
			t, e := NewTranslator(o)
			if e != nil {
				return
			}
			c, e := t.Translate(v)
			if e != nil {
				return
			}
			k.Clause = c
		}
		return
	default:
		err = fmt.Errorf("can't use datasource type=%v as dateset's datasource", source.Type)
		return
	}
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
