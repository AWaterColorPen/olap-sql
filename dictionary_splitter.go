package olapsql

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/models"
	"github.com/awatercolorpen/olap-sql/api/types"
)

func getSourceFromTranslator(translator Translator) *models.DataSource {
	current := translator.GetCurrent()
	adapter := translator.GetAdapter()
	source, _ := adapter.GetSourceByKey(current)
	return source
}

type normalClauseSplitter struct {
	Translator Translator

	CandidateList []*types.DataSource
	Candidate     map[string]*types.DataSource
	SplitQuery    map[string]*types.Query

	Clause *types.NormalClause
	Query  *types.Query
	DBType types.DBType
}

func NewNormalClauseSplitter(translator Translator, clause *types.NormalClause, query *types.Query, dbType types.DBType) (*normalClauseSplitter, error) {
	adapter := translator.GetAdapter()
	source := getSourceFromTranslator(translator)

	candidateList := []*types.DataSource{convertDataSourceToDataSource(source)}
	for _, hit := range source.GetGetDependencyKey() {
		in, _ := adapter.GetSourceByKey(hit)
		candidateList = append(candidateList, convertDataSourceToDataSource(in))
	}

	candidate := map[string]*types.DataSource{}
	splitQuery := map[string]*types.Query{}
	for _, v := range candidateList {
		candidate[v.Name] = v
		splitQuery[v.Name] = &types.Query{}
	}
	splitter := &normalClauseSplitter{
		Translator:    translator,
		CandidateList: candidateList,
		Candidate:     candidate,
		SplitQuery:    splitQuery,
		Clause:        clause,
		Query:         query,
		DBType:        dbType,
	}
	return splitter, nil
}

func (n *normalClauseSplitter) Run() error {
	source := getSourceFromTranslator(n.Translator)
	switch source.Type {
	case types.DataSourceTypeFact:
		return n.factRun()
	case types.DataSourceTypeFactDimensionJoin:
		return n.joinRun()
	case types.DataSourceTypeMergeJoin:
		return n.mergeRun()
	default:
		return fmt.Errorf("can't use datasource type=%v as dateset's datasource", source.Type)
	}
}

func (n *normalClauseSplitter) factRun() error {
	source := n.GetSelfCandidate()
	n.Clause.DataSource = append(n.Clause.DataSource, source)
	return nil
}

func (n *normalClauseSplitter) joinRun() error {
	candidate := n.GetOtherCandidate()
	for _, v := range candidate {
		n.Clause.DataSource = append(n.Clause.DataSource, v)
	}
	joins, err := n.buildJoins()
	if err != nil {
		return err
	}
	n.Clause.Joins = joins
	return nil
}

func (n *normalClauseSplitter) mergeRun() error {
	if err := n.joinRun(); err != nil {
		return err
	}
	if err := n.split(); err != nil {
		return err
	}
	oq := n.GetOtherSplitQuery()
	for k, v := range oq {
		option := &TranslatorOption{
			Adapter: n.Translator.GetAdapter(),
			Query:   v,
			DBType:  n.DBType,
			Current: k,
		}
		translator, err := NewTranslator(option)
		if err != nil {
			return err
		}
		clause, err := translator.Translate(v)
		if err != nil {
			return err
		}
		n.Candidate[k].Clause = clause
	}

	sq := n.GetSelfSplitQuery()
	n.Clause.Filters = sq.Filters
	return nil
}

func (n *normalClauseSplitter) Polish() map[string]*types.Query {
	n.polish()
	return n.SplitQuery
}

func (n *normalClauseSplitter) GetSelfSplitQuery() *types.Query {
	current := n.Translator.GetCurrent()
	return n.Polish()[current]
}

func (n *normalClauseSplitter) GetOtherSplitQuery() map[string]*types.Query {
	other := map[string]*types.Query{}
	current := n.Translator.GetCurrent()
	for k, v := range n.Polish() {
		if k != current {
			other[k] = v
		}
	}
	return other
}

func (n *normalClauseSplitter) GetSelfCandidate() *types.DataSource {
	current := n.Translator.GetCurrent()
	return n.Candidate[current]
}

func (n *normalClauseSplitter) GetOtherCandidate() []*types.DataSource {
	var other []*types.DataSource
	current := n.Translator.GetCurrent()
	for _, v := range n.CandidateList {
		if v.Name != current {
			other = append(other, v)
		}
	}
	return other
}

func (n *normalClauseSplitter) split() error {
	for _, m := range n.Clause.Metrics {
		kv, err := n.splitMetric(m)
		if err != nil {
			return err
		}
		n.addMetric(kv)
	}
	for _, d := range n.Clause.Dimensions {
		kv, err := n.splitDimension(d)
		if err != nil {
			return err
		}
		n.addDimension(kv)
	}
	for _, f := range n.Query.Filters {
		kv, err := n.splitFilter(f)
		if err != nil {
			return err
		}
		n.addFilter(kv)
	}
	return nil
}

func (n *normalClauseSplitter) splitMetric(metric *types.Metric) (map[string][]string, error) {
	out := map[string][]string{}
	if len(metric.Table) > 0 {
		out[metric.Table] = append(out[metric.Table], metric.Name)
	}
	for _, child := range metric.Children {
		kv, err := n.splitMetric(child)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k] = append(out[k], v...)
		}
	}
	return out, nil
}

func (n *normalClauseSplitter) splitDimension(dimension *types.Dimension) (map[string][]string, error) {
	out := map[string][]string{}
	if len(dimension.Table) > 0 {
		out[dimension.Table] = append(out[dimension.Table], dimension.Name)
	}
	for _, child := range dimension.Dependency {
		kv, err := n.splitDimension(child)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k] = append(out[k], v...)
		}
	}
	return out, nil
}

func (n *normalClauseSplitter) splitFilter(filter *types.Filter) (map[string][]*types.Filter, error) {
	out := map[string][]*types.Filter{}
	target, err := getOneFilterTargetTable(n.Translator, filter)
	if err != nil {
		return nil, err
	}
	output, err := convertOneFilterToTargetTables(n.Translator, filter, target)
	if err != nil {
		return nil, err
	}
	for k, v := range output {
		out[k] = append(out[k], v)
	}
	return out, nil
}

func (n *normalClauseSplitter) polish() {
	for _, v := range n.SplitQuery {
		linq.From(v.Metrics).Distinct().ToSlice(&v.Metrics)
		linq.From(v.Dimensions).Distinct().ToSlice(&v.Dimensions)
	}
}

func (n *normalClauseSplitter) addMetric(metric map[string][]string) {
	for k, v := range metric {
		n.SplitQuery[k].Metrics = append(n.SplitQuery[k].Metrics, v...)
	}
}

func (n *normalClauseSplitter) addDimension(dimension map[string][]string) {
	for k, v := range dimension {
		n.SplitQuery[k].Dimensions = append(n.SplitQuery[k].Dimensions, v...)
	}
}

func (n *normalClauseSplitter) addFilter(filter map[string][]*types.Filter) {
	for k, v := range filter {
		n.SplitQuery[k].Filters = append(n.SplitQuery[k].Filters, v...)
	}
}

func buildHitMergeJoinDimension(translator Translator, query *types.Query) (map[string]bool, error) {
	current := translator.GetCurrent()
	adapter := translator.GetAdapter()
	source, _ := adapter.GetSourceByKey(current)
	if source.Type != types.DataSourceTypeMergeJoin {
		return nil, fmt.Errorf("the source is not a merged join data source")
	}

	set := map[string]bool{}
	for _, v := range query.Dimensions {
		set[v] = true
	}
	hit := map[string]bool{}
	for _, v := range source.MergeJoin[0].Dimension {
		if _, ok := set[v]; ok {
			hit[v] = true
		}
	}
	return hit, nil
}

func (n *normalClauseSplitter) buildDimensionJoin() []*types.Join {
	candidate := n.Candidate
	dGraph := n.Translator.GetDependencyGraph()
	source := getSourceFromTranslator(n.Translator)
	var joins []*types.Join
	for _, v := range source.DimensionJoin {
		s1, ok1 := candidate[v.Get1().DataSource]
		s2, ok2 := candidate[v.Get2().DataSource]
		if !ok1 || !ok2 {
			continue
		}
		ds1, dl1, ds2, dl2 := v.Get1().DataSource, v.Get1().Dimension, v.Get2().DataSource, v.Get2().Dimension
		var on []*types.JoinOn
		for i := 0; i < len(dl1); i++ {
			k1 := fmt.Sprintf("%v.%v", ds1, dl1[i])
			k2 := fmt.Sprintf("%v.%v", ds2, dl2[i])
			d1, _ := dGraph.GetDimension(k1)
			d2, _ := dGraph.GetDimension(k2)
			key1, _ := d1.Expression()
			key2, _ := d2.Expression()
			on = append(on, &types.JoinOn{Key1: models.GetNameFromKey(key1), Key2: models.GetNameFromKey(key2)})
		}

		j := &types.Join{DataSource1: s1, DataSource2: s2, On: on}
		joins = append(joins, j)
	}
	return joins
}

func (n *normalClauseSplitter) buildMergeJoin() []*types.Join {
	candidate := n.Candidate
	dGraph := n.Translator.GetDependencyGraph()
	source := getSourceFromTranslator(n.Translator)
	hitDimension, _ := buildHitMergeJoinDimension(n.Translator, n.Query)
	var joins []*types.Join
	for i := 2; i < len(source.MergeJoin); i++ {
		s1, ok1 := candidate[source.MergeJoin[1].DataSource]
		s2, ok2 := candidate[source.MergeJoin[i].DataSource]
		if !ok1 || !ok2 {
			continue
		}
		ds1, dl1 := source.MergeJoin[1].DataSource, source.MergeJoin[1].Dimension
		ds2, dl2 := source.MergeJoin[i].DataSource, source.MergeJoin[i].Dimension
		var on []*types.JoinOn
		for j := 0; j < len(dl1); j++ {
			if _, ok := hitDimension[source.MergeJoin[0].Dimension[j]]; !ok {
				continue
			}
			k1 := fmt.Sprintf("%v.%v", ds1, dl1[j])
			k2 := fmt.Sprintf("%v.%v", ds2, dl2[j])
			d1, _ := dGraph.GetDimension(k1)
			d2, _ := dGraph.GetDimension(k2)
			key1, _ := d1.Alias()
			key2, _ := d2.Alias()
			on = append(on, &types.JoinOn{Key1: models.GetNameFromKey(key1), Key2: models.GetNameFromKey(key2)})
		}
		j := &types.Join{DataSource1: s1, DataSource2: s2, On: on}
		joins = append(joins, j)
	}
	return joins
}

func (n *normalClauseSplitter) buildJoins() ([]*types.Join, error) {
	var joins []*types.Join
	joins = append(joins, n.buildDimensionJoin()...)
	joins = append(joins, n.buildMergeJoin()...)
	return joins, nil
}

func getOneFilterTargetTable(translator Translator, in *types.Filter) ([]string, error) {
	out := &types.Filter{
		OperatorType: in.OperatorType,
		Value:        in.Value,
	}
	if !out.OperatorType.IsTree() {
		c, err := getColumn(translator, translator.GetCurrent(), in.Name)
		if err != nil {
			return nil, err
		}
		return c.GetTables(), nil
	}

	var cTable [][]string
	for _, v := range in.Children {
		table, err := getOneFilterTargetTable(translator, v)
		if err != nil {
			return nil, err
		}
		cTable = append(cTable, table)
	}

	for i := 1; i < len(cTable); i++ {
		if err := isSameColumnTables(cTable[i-1], cTable[i]); err != nil {
			return nil, err
		}
	}
	return cTable[0], nil
}

func convertOneFilterToTargetTable(translator Translator, in *types.Filter, target string) error {
	if target == translator.GetCurrent() {
		return nil
	}
	if !in.OperatorType.IsTree() {
		c, err := getColumn(translator, translator.GetCurrent(), in.Name)
		if err != nil {
			return err
		}
		in.Name = c.GetTargetFieldName(target)
		return nil
	}

	for _, v := range in.Children {
		err := convertOneFilterToTargetTable(translator, v, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertOneFilterToTargetTables(translator Translator, in *types.Filter, target []string) (map[string]*types.Filter, error) {
	out := map[string]*types.Filter{}
	for _, v := range target {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		input := &types.Filter{}
		if err = json.Unmarshal(b, input); err != nil {
			return nil, err
		}
		if err = convertOneFilterToTargetTable(translator, input, v); err != nil {
			return nil, err
		}
		out[v] = input
	}
	return out, nil
}

func convertDataSourceToDataSource(in *models.DataSource) *types.DataSource {
	return &types.DataSource{Database: in.Database, Name: in.Name, AliasName: in.Alias, Type: in.Type}
}
