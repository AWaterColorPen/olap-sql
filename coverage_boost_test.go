package olapsql_test

// coverage_boost_test.go — Phase 3 coverage boost
// Targets: Clients.SetLogger, Clients.BuildSQL, Manager.SetLogger, Manager.BuildSQL,
//          NewTranslator (direct-SQL path), FileAdapter.GetMetricsBySource/GetDimensionsBySource,
//          database.getDialect unsupported-type error path.

import (
	"testing"

	olapsql "github.com/awatercolorpen/olap-sql"
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

// ---------------------------------------------------------------------------
// Clients.SetLogger
// ---------------------------------------------------------------------------

func TestClients_SetLogger(t *testing.T) {
	clients, err := newClients(t.TempDir())
	assert.NoError(t, err)

	// SetLogger should not panic; pass a no-op logger.
	clients.SetLogger(logger.Discard)
}

// ---------------------------------------------------------------------------
// Clients.BuildSQL
// ---------------------------------------------------------------------------

func TestClients_BuildSQL(t *testing.T) {
	// BuildSQL uses DryRun mode, so no data loading is required.
	m, err := newManager(t)
	assert.NoError(t, err)

	query := MockQuery1()

	dictionary, err := m.GetDictionary()
	assert.NoError(t, err)

	clause, err := dictionary.Translate(query)
	assert.NoError(t, err)

	clients, err := m.GetClients()
	assert.NoError(t, err)

	sqlStr, err := clients.BuildSQL(clause)
	assert.NoError(t, err)
	assert.NotEmpty(t, sqlStr)
}

// ---------------------------------------------------------------------------
// Manager.SetLogger
// ---------------------------------------------------------------------------

func TestManager_SetLogger(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)

	// Should silently succeed (logs go to Discard).
	m.SetLogger(logger.Discard)
}

// ---------------------------------------------------------------------------
// Manager.BuildSQL
// ---------------------------------------------------------------------------

func TestManager_BuildSQL(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))

	query := MockQuery1()
	sql, err := m.BuildSQL(query)
	assert.NoError(t, err)
	assert.NotEmpty(t, sql)
}

// ---------------------------------------------------------------------------
// NewTranslator — direct-SQL path (query.Sql != "")
// ---------------------------------------------------------------------------

func TestNewTranslator_DirectSQL(t *testing.T) {
	dict, err := newDictionary()
	assert.NoError(t, err)

	query := &types.Query{
		DataSetName: mockWikiStatDataSet,
		Sql:         "SELECT 1",
	}

	clause, err := dict.Translate(query)
	assert.NoError(t, err)
	assert.NotNil(t, clause)

	// Verify the direct-SQL path preserved the raw SQL string.
	sqlClause, ok := clause.(*types.SqlClause)
	assert.True(t, ok, "expected clause to be *types.SqlClause")
	assert.Equal(t, "SELECT 1", sqlClause.Sql)
}

// ---------------------------------------------------------------------------
// FileAdapter.GetMetricsBySource / GetDimensionsBySource
// ---------------------------------------------------------------------------

func TestFileAdapter_GetMetricsBySource(t *testing.T) {
	adapter, err := newFileAdapter()
	assert.NoError(t, err)

	// "wikistat" is a known source in the test dictionary.
	metrics := adapter.GetMetricsBySource("wikistat")
	assert.NotEmpty(t, metrics, "expected at least one metric for source 'wikistat'")

	// Non-existent source should return empty slice (not an error).
	none := adapter.GetMetricsBySource("__nonexistent__")
	assert.Empty(t, none)
}

func TestFileAdapter_GetDimensionsBySource(t *testing.T) {
	adapter, err := newFileAdapter()
	assert.NoError(t, err)

	dims := adapter.GetDimensionsBySource("wikistat")
	assert.NotEmpty(t, dims, "expected at least one dimension for source 'wikistat'")

	none := adapter.GetDimensionsBySource("__nonexistent__")
	assert.Empty(t, none)
}

// ---------------------------------------------------------------------------
// DBOption.NewDB — unsupported type returns error
// ---------------------------------------------------------------------------

func TestDBOption_NewDB_UnsupportedType(t *testing.T) {
	opt := &olapsql.DBOption{
		DSN:  "whatever",
		Type: types.DBType("__unknown__"),
	}
	_, err := opt.NewDB()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported db type")
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newDictionary() (*olapsql.Dictionary, error) {
	do := getDictionaryOption()
	return olapsql.NewDictionary(do)
}

func newFileAdapter() (olapsql.IAdapter, error) {
	opt := &olapsql.AdapterOption{
		Type: olapsql.FILEAdapter,
		Dsn:  "test/dictionary.sqlite.toml",
	}
	return olapsql.NewAdapter(opt)
}
