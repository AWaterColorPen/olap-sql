package olapsql

import (
	"fmt"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBType string

const (
	DBTypeSQLite     DBType = "sqlite"
	DBTypeMySQL      DBType = "mysql"
	DBTypePostgre    DBType = "postgres"
	DBTypeClickHouse DBType = "clickhouse"
)

type DBOption struct {
	Debug bool   `json:"debug"`
	DSN   string `json:"dsn"`
	Type  DBType `json:"type"`
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

func getDialect(ty DBType, dsn string) (gorm.Dialector, error) {
	switch ty {
	case DBTypeSQLite:
		return sqlite.Open(dsn), nil
	case DBTypeMySQL:
		return mysql.Open(dsn), nil
	case DBTypePostgre:
		return postgres.Open(dsn), nil
	case DBTypeClickHouse:
		return clickhouse.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %v", ty)
	}
}
