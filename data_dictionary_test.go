package olapsql_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewDataDictionary(t *testing.T) {
	_, err := newDataDictionary(t.TempDir())
	assert.NoError(t, err)
}

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func newDataDictionary(sqlitePath string) (*olapsql.DataDictionary, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	option := &olapsql.DataDictionaryOption{DB: db}
	return olapsql.NewDataDictionary(option)
}