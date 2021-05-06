package olapsql_test

import (
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestNewClients(t *testing.T) {
	_, err := newClients(t.TempDir())
	assert.NoError(t, err)
}

func TestClients_Translator(t *testing.T) {
	d, err := newDataDictionary(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, dataDictionaryMockData(d))

	query := &types.Query{
		Metrics:    []string{"cost", "click"},
		Dimensions: []string{"account_id", "account_name"},
		Filters: []*types.Filter{
			{OperatorType: types.FilterOperatorTypeIn, Name: "partition_time", Value: []interface{}{"2021-04-29", "2021-04-30"}},
		},
		DataSet: mockDataSet,
	}

	translator, err := d.Translator(query)
	assert.NoError(t, err)
	request, err := translator.Translate(query)
	assert.NoError(t, err)

	client, err := newClients(t.TempDir())
	assert.NoError(t, err)
	db, err := client.Get(request.(*types.Request).DataSource.Type, "")
	assert.NoError(t, err)
	sql, err := request.(*types.Request).Clause(db)
	assert.NoError(t, err)
	t.Log(sql.Statement.SQL.String())
	_, _ = sql.Debug().Rows()
	// t.Log(statement.Vars)
}

func newClients(sqlitePath string) (olapsql.Clients, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	client := olapsql.Clients{}
	client.RegisterByKV(types.DataSourceTypeClickHouse, mockDataSet, db)
	client.RegisterByKV(types.DataSourceTypeClickHouse, "", db)
	return client, nil
}
