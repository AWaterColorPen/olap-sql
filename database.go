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

type DBOption struct {
	Debug bool         `json:"debug"`
	DSN   string       `json:"dsn"`
	Type  types.DBType `json:"type"`
}

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
