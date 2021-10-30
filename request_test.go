package olapsql_test

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))
	request := Request1()
	client, err := m.GetClients()
	assert.NoError(t, err)
	sql, err := client.BuildSQL(request)
	fmt.Println(sql)
}

func Request1() *types.NormalClause {
	request := &types.NormalClause{
		DBType:  types.DBTypeSQLite,
		Dataset: mockWikiStatDataSet,
		Metrics: []*types.Metric{
			{
				Type:      types.MetricTypeValue,
				Table:     "t1",
				Name:      "hits",
				FieldName: "hits",
				DBType:    types.DBTypeSQLite,
			},
		},
		DataSource: &types.DataSource{
			Alias: "t1",
			Clause: &types.NormalClause{
				DBType:  types.DBTypeSQLite,
				Dataset: mockWikiStatDataSet,
				Metrics: []*types.Metric{
					{
						Type:      types.MetricTypeValue,
						Table:     "wikistat",
						Name:      "hits",
						FieldName: "hits",
						DBType:    types.DBTypeSQLite,
					},
				},
				DataSource: &types.DataSource{
					Name: "wikistat",
				},
			},
		},
		Joins: []*types.Join{
			{
				DataSource1: &types.DataSource{
					Alias: "t1",
					Clause: &types.NormalClause{
						DBType:  types.DBTypeSQLite,
						Dataset: mockWikiStatDataSet,
						Metrics: []*types.Metric{
							{
								Type:      types.MetricTypeValue,
								Table:     "wikistat",
								Name:      "hits",
								FieldName: "hits",
								DBType:    types.DBTypeSQLite,
							},
						},
						DataSource: &types.DataSource{
							Name: "wikistat",
						},
					},
				},
				DataSource2: &types.DataSource{
					Alias: "t2",
					Clause: &types.NormalClause{
						DBType:  types.DBTypeSQLite,
						Dataset: mockWikiStatDataSet,
						Metrics: []*types.Metric{
							{
								Type:      types.MetricTypeValue,
								Table:     "wikistat",
								Name:      "hits",
								FieldName: "hits",
								DBType:    types.DBTypeSQLite,
							},
						},
						DataSource: &types.DataSource{
							Name: "wikistat",
						},
					},
				},
				On: []*types.JoinOn{
					{
						Key1: "hits",
						Key2: "hits",
					},
				},
			},
		},
		Filters: []*types.Filter{
			{
				OperatorType: types.FilterOperatorTypeGreaterEquals,
				ValueType:    types.ValueTypeInteger,
				Name:         "hits",
				Value:        []interface{}{1},
			},
		},
	}
	return request
}
