package types

type DBType string

const (
	DBTypeSQLite     DBType = "sqlite"
	DBTypeMySQL      DBType = "mysql"
	DBTypePostgre    DBType = "postgres"
	DBTypeClickHouse DBType = "clickhouse"
	DBTypeDruid      DBType = "druid"
	DBTypeKylin      DBType = "kylin"
)
