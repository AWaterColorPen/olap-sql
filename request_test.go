package olapsql_test

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/sirupsen/logrus"
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
	db, err := client.BuildDB(request)
	logrus.Info(request.BuildSQL(db))
}


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