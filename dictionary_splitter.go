package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Splitter interface {
}

type normalClauseSplitter struct {
	Candidate map[string]*types.DataSource
	SplitQuery map[string]*types.Query
	Clause *types.NormalClause
}

func NewNormalClauseSplitter(clause *types.NormalClause, source []*types.DataSource) (*normalClauseSplitter, error) {
	candidate := map[string]*types.DataSource{}
	splitQuery := map[string]*types.Query{}
	for _, v := range source {
		candidate[v.Name] = v
		splitQuery[v.Name] = &types.Query{}
	}
	splitter := &normalClauseSplitter{
		Candidate: candidate,
		SplitQuery: splitQuery,
		Clause:     clause,
	}
	return splitter, nil
}

func (n *normalClauseSplitter) Split() (map[*types.DataSource]*types.Query, error) {
	out := map[*types.DataSource]*types.Query{}
	for _, v := range n.Candidate {
		out[v] = &types.Query{}
	}

	for _, m := range n.Clause.Metrics {
		kv, err := n.splitMetric(m)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k].Metrics = append(out[k].Metrics, v...)
		}
	}

	for _, d := range n.Clause.Dimensions {
		kv, err := n.splitDimension(d)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k].Dimensions = append(out[k].Dimensions, v...)
		}
	}
	for _, f := range n.Clause.Filters {
		_, err := n.splitFilter(f)
		if err != nil {
			return nil, err
		}
		// for k, v := range kv {
		// 	// out[k].Filters = v
		// }
	}
	return out, nil
}

func (n *normalClauseSplitter) splitMetric(metric *types.Metric) (map[*types.DataSource][]string, error) {
	out := map[*types.DataSource][]string{}
	if len(metric.Table) > 0 {
		if hit, ok := n.Candidate[metric.Table]; ok {
			out[hit] = append(out[hit], metric.Name)
		}
	}
	for _, child := range metric.Children {
		kv, err := n.splitMetric(child)
		if err != nil {
			return nil, err
		}
		mergeSplitMap(kv, &out)
	}
	return out, nil
}

func (n *normalClauseSplitter) splitDimension(dimension *types.Dimension) (map[*types.DataSource][]string, error) {
	out := map[*types.DataSource][]string{}
	if len(dimension.Table) > 0 {
		if hit, ok := n.Candidate[dimension.Table]; ok {
			out[hit] = append(out[hit], dimension.Name)
		}
	}
	for _, child := range dimension.Dependency {
		kv, err := n.splitDimension(child)
		if err != nil {
			return nil, err
		}
		mergeSplitMap(kv, &out)
	}
	return out, nil
}

func (n *normalClauseSplitter) splitFilter(filter *types.Filter) (map[*types.DataSource][]string, error) {
	out := map[*types.DataSource][]string{}
	if filter == nil {
		return out, nil
	}

	if len(filter.Table) > 0 {
		if hit, ok := n.Candidate[filter.Table]; ok {
			out[hit] = append(out[hit], filter.Name)
		}
	}
	for _, child := range filter.Children {
		kv, err := n.splitFilter(child)
		if err != nil {
			return nil, err
		}
		mergeSplitMap(kv, &out)
	}
	return out, nil
}

func mergeSplitMap(from map[*types.DataSource][]string, to *map[*types.DataSource][]string) {
	for k, v := range from {
		(*to)[k] = append((*to)[k], v...)
	}
}
