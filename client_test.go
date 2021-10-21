package olapsql_test

import (
	"gorm.io/gorm"
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
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

	dictionary, err := m.GetDictionary()
	assert.NoError(t, err)
	request, err := dictionary.Translate(query)
	assert.NoError(t, err)

	client, err := m.GetClients()
	assert.NoError(t, err)
	db, err := client.SubmitClause(request)

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

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
