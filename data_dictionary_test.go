package olapsql_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/olap-sql/dictionary"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	mockDataSet = "mock-dataset"
)

func TestNewDataDictionary(t *testing.T) {
	d, err := newDataDictionary(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, MockWikiStatDataDictionary(d))
}

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func newDataDictionary(sqlitePath string) (*dictionary.Dictionary, error) {
	option := &dictionary.DictionaryOption{
		dictionary.AdapterOption{
			Type: dictionary.FILEadapter,
			Dsn: "filetest/test.yaml",
		},
	}
	return dictionary.NewDictionary(option)
}
