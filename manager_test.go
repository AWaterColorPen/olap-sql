package olapsql_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	_, err = m.GetDataDictionary()
	assert.NoError(t, err)
	_, err = m.GetClients()
	assert.NoError(t, err)
}

func TestManager_RunChan(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))

	query := MockQuery1()
	result, err := m.RunChan(query)
	assert.NoError(t, err)
	MockQuery1ResultAssert(t, result)
}

func newManager(tb testing.TB) (*olapsql.Manager, error) {
	configuration := &olapsql.Configuration{
		ClientsOption: map[string]*olapsql.DBOption{
			string(types.DataSourceTypeClickHouse): getOlapDBOption(tb),
		},
		DataDictionaryOption: &olapsql.DataDictionaryOption{
			DBOption: olapsql.DBOption{DSN: filepath.Join(tb.TempDir(), "sqlite"), Type: olapsql.DBTypeSQLite},
		},
	}

	return olapsql.NewManager(configuration)
}

func getOlapDBOption(tb testing.TB) *olapsql.DBOption {
	if DataWithClickhouse() {
		return &olapsql.DBOption{DSN: "tcp://localhost:9000?database=default", Type: olapsql.DBTypeClickHouse}
	}
	return &olapsql.DBOption{DSN: filepath.Join(tb.TempDir(), "sqlite"), Type: olapsql.DBTypeSQLite}
}
