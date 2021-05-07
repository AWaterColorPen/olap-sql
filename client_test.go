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

func TestClients_SubmitClause(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))

	query := MockQuery1()

	dictionary, err := m.GetDataDictionary()
	assert.NoError(t, err)
	translator, err := dictionary.Translator(query)
	assert.NoError(t, err)
	request, err := translator.Translate(query)
	assert.NoError(t, err)

	client, err := m.GetClients()
	assert.NoError(t, err)
	db, err := client.SubmitClause(request.(*types.Request))
	assert.NoError(t, err)

	results, err := olapsql.RunSync(db.Debug())
	assert.NoError(t, err)
	table, err := olapsql.BuildTableResultSync(query, results)
	assert.NoError(t, err)
	assert.Len(t, table.Header, 7)
	assert.Equal(t,"date", table.Header[0])
	assert.Equal(t,"source_avg", table.Header[6])
	assert.Len(t, table.Rows, 3)
	assert.Len(t, table.Rows[0], 7)
	assert.Equal(t, float64(10244), table.Rows[0][3])
	assert.Equal(t, 0.013861772745021476, table.Rows[0][5])
	assert.Equal(t, 2.52971, table.Rows[0][6])
}

func newClients(sqlitePath string) (olapsql.Clients, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	client := olapsql.Clients{}
	client.RegisterByKV(types.DataSourceTypeClickHouse, mockWikiStatDataSet, db)
	client.RegisterByKV(types.DataSourceTypeClickHouse, "", db)
	return client, nil
}
