package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type Clause interface {
	GetDBType() types.DBType
	GetDataset() string
	BuildDB(tx *gorm.DB) (*gorm.DB, error)
	BuildSql(tx *gorm.DB) (string, error)
}
