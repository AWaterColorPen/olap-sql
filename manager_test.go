package olapsql_test

import (
	"path/filepath"
	"testing"

	olapsql "github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/awatercolorpen/olap-sql/dictionary"
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
	testMockQuery(t, MockQuery1(), MockQuery1ResultAssert)
}

func newManager(tb testing.TB) (*olapsql.Manager, error) {
	k, v := getOlapDBOption(tb)
	configuration := &olapsql.Configuration{
		ClientsOption: map[string]*olapsql.DBOption{
			k: v,
		},
		DataDictionaryOption: &dictionary.DictionaryOption{
			AdapterOption: dictionary.AdapterOption{
				Type: dictionary.FILEadapter,
				Dsn:  "filetest/test.yaml",
			},
		},
	}

	return olapsql.NewManager(configuration)
}

func getOlapDBOption(tb testing.TB) (string, *olapsql.DBOption) {
	if DataWithClickhouse() {
		return string(types.DataSourceTypeClickHouse), &olapsql.DBOption{DSN: "tcp://localhost:9000?database=default", Type: olapsql.DBTypeClickHouse}
	}
	return string(types.DataSourceTypeUnknown), &olapsql.DBOption{DSN: filepath.Join(tb.TempDir(), "sqlite"), Type: olapsql.DBTypeSQLite, Debug: true}
}
