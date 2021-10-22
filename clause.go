package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type Clause interface {
	GetDBType() types.DBType
	GetDataset() string
	BuildDB(tx *gorm.DB) (*gorm.DB, error)
	BuildSQL(tx *gorm.DB) (string, error)
}

type baseClause struct {
	DBType  types.DBType
	Dataset string
}

func (b *baseClause) GetDBType() types.DBType {
	return b.DBType
}

func (b *baseClause) GetDataset() string {
	return b.Dataset
}

type sqlClause struct {
	baseClause
	sql string
}

func (s *sqlClause) BuildDB(tx *gorm.DB) (*gorm.DB, error)  {
	return tx.Raw(s.sql), nil
}

func (s *sqlClause) BuildSQL(*gorm.DB) (string, error)  {
	return s.sql, nil
}
