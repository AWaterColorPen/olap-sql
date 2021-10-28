package olapsql_test

import (
	"path/filepath"
	"testing"

	olapsql "github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	_, err = m.GetDictionary()
	assert.NoError(t, err)
	_, err = m.GetClients()
	assert.NoError(t, err)
}

func TestManager_RunChan(t *testing.T) {
	testMockQuery(t, MockQuery1(), MockQuery1ResultAssert)
}

func newManager(tb testing.TB) (*olapsql.Manager, error) {
	k, v := getOlapDBOption(tb)
	do := getDictionaryOption()
	configuration := &olapsql.Configuration{
		ClientsOption:    map[string]*olapsql.DBOption{k: v},
		DictionaryOption: do,
	}

	return olapsql.NewManager(configuration)
}

func getOlapDBOption(tb testing.TB) (string, *olapsql.DBOption) {
	if DataWithClickhouse() {
		return string(types.DBTypeClickHouse), &olapsql.DBOption{DSN: "tcp://localhost:9000?database=default", Type: types.DBTypeClickHouse}
	}
	return string(types.DBTypeSQLite), &olapsql.DBOption{DSN: filepath.Join(tb.TempDir(), "sqlite"), Type: types.DBTypeSQLite, Debug: true}
}

func getDictionaryOption() *olapsql.Option {
	dsn := "test/dictionary.sqlite.toml"
	if DataWithClickhouse() {
		dsn = "test/dictionary.ck.toml"
	}
	return &olapsql.Option{
		AdapterOption: olapsql.AdapterOption{
			Type: olapsql.FILEAdapter,
			Dsn:  dsn,
		},
	}
}
