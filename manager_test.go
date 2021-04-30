package olapsql_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	configuration := &olapsql.Configuration{
		ClientsOption: map[string]*olapsql.DBOption{
			string(types.DataSourceTypeClickHouse): {DSN: filepath.Join(t.TempDir(), "sqlite"), Type: olapsql.DBTypeSQLite},
		},
		DataDictionaryOption: &olapsql.DataDictionaryOption{
			DBOption: olapsql.DBOption{DSN: filepath.Join(t.TempDir(), "sqlite"), Type: olapsql.DBTypeSQLite},
		},
	}
	m, err := olapsql.NewManager(configuration)
	assert.NoError(t, err)
	dictionary, err := m.GetDataDictionary()
	assert.NoError(t, err)
	assert.NoError(t, dataDictionaryMockData(dictionary))
	clients, err := m.GetClients()
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
}
