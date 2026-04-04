package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBOption holds the connection parameters for a single database instance.
type DBOption struct {
	// Debug enables GORM's debug mode, which prints every generated SQL statement.
	Debug bool `json:"debug"`

	// DSN is the data source name (connection string) for the database.
	// Format depends on the driver, e.g.:
	//   ClickHouse: "clickhouse://user:pass@host:9000/db"
	//   MySQL:      "user:pass@tcp(host:3306)/db?charset=utf8"
	//   PostgreSQL: "host=host user=user password=pass dbname=db port=5432"
	//   SQLite:     "/path/to/file.db"
	DSN string `json:"dsn"`

	// Type identifies the database engine. Supported values: clickhouse, mysql, postgre, sqlite.
	Type types.DBType `json:"type"`
}

// NewDB opens a database connection using the DBOption settings.
// Returns a configured *gorm.DB or an error if the DSN or type is invalid.
func (o *DBOption) NewDB() (*gorm.DB, error) {
	dialect, err := getDialect(o.Type, o.DSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if o.Debug {
		db = db.Debug()
	}
	return db, nil
}

// getDialect maps a DBType to the corresponding GORM Dialector.
func getDialect(ty types.DBType, dsn string) (gorm.Dialector, error) {
	switch ty {
	case types.DBTypeSQLite:
		return sqlite.Open(dsn), nil
	case types.DBTypeMySQL:
		return mysql.Open(dsn), nil
	case types.DBTypePostgre:
		return postgres.Open(dsn), nil
	case types.DBTypeClickHouse:
		return clickhouse.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %v", ty)
	}
}
