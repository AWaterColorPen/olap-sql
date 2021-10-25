package olapsql_test

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

func Request1() *types.Request {
	request := &types.Request{
		DBType: types.DBTypeSQLite,
		Dataset: mockWikiStatDataSet,
		Metrics: []*types.Metric{
			{
				Type:      types.MetricTypeValue,
				Table:     "wikistat",
				Name:      "hits",
				FieldName: "hits",
				DBType:    types.DBType(types.DataSourceTypeUnknown),
			},
		},
	}
	return request
}