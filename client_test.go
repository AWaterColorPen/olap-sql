package olapsql_test

import (
	"testing"

	olapsql "github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestNewClients(t *testing.T) {
	_, err := newClients(t.TempDir())
	assert.NoError(t, err)
}

func TestClients_SubmitClause(t *testing.T) {
	assert.NoError(t, MockWikiStatDataToJson())
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))

	query := MockQuery1()

	dictionary, err := m.GetDataDictionary()
	assert.NoError(t, err)
	request, err := dictionary.Translate(query)
	assert.NoError(t, err)

	client, err := m.GetClients()
	assert.NoError(t, err)
	db, err := client.SubmitClause(request)
	assert.NoError(t, err)

	rows, err := olapsql.RunSync(db)
	assert.NoError(t, err)
	result, err := olapsql.BuildResultSync(query, rows)
	assert.NoError(t, err)
	MockQuery1ResultAssert(t, result)
}

func newClients(sqlitePath string) (olapsql.Clients, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	client := olapsql.Clients{}
	client.RegisterByKV(types.DataSourceTypeUnknown, mockWikiStatDataSet, db)
	client.RegisterByKV(types.DataSourceTypeUnknown, "", db)
	return client, nil
}
