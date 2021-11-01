package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Splitter interface {
}

type normalClauseSplitter struct {
	Candidate map[string]*types.DataSource
}

func NewNormalClauseSplitter(source []*types.DataSource) (*normalClauseSplitter, error) {
	candidate := map[string]*types.DataSource{}
	for _, v := range candidate {
		candidate[v.Name] = v
	}
	splitter := &normalClauseSplitter{
		Candidate: candidate,
	}
	return splitter, nil
}

func (n *normalClauseSplitter) Split(clause *types.NormalClause) (map[*types.DataSource]*types.Query, error) {
	out := map[*types.DataSource]*types.Query{}
	for _, v := range n.Candidate {
		out[v] = &types.Query{}
	}

	for _, m := range clause.Metrics {
		kv, err := n.splitMetric(m)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k].Metrics = v
		}
	}

	for _, d := range clause.Dimensions {
		kv, err := n.splitDimension(d)
		if err != nil {
			return nil, err
		}
		for k, v := range kv {
			out[k].Dimensions = v
		}
	}
	for _, f := range clause.Filters {
		_, err := n.splitFilter(f)
		if err != nil {
			return nil, err
		}
		// for k, v := range kv {
		// 	out[k].Filters = v
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

	kv, err := n.splitFilter(metric.Filter)
	if err != nil {
		return nil, err
	}
	mergeSplitMap(kv, &out)
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

func (n *normalClauseSplitter) splitDimension(dimension *types.Dimension) (map[*types.DataSource][]string, error) {
	out := map[*types.DataSource][]string{}
	if hit, ok := n.Candidate[dimension.Table]; ok {
		out[hit] = append(out[hit], dimension.Name)
	}
	return out, nil
}

func (n *normalClauseSplitter) splitOrderBy(_ *types.OrderBy) (map[*types.DataSource][]string, error) {
	return nil, nil
}

func mergeSplitMap(from map[*types.DataSource][]string, to *map[*types.DataSource][]string) {
	for k, v := range from {
		(*to)[k] = append((*to)[k], v...)
	}
}
