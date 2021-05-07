package olapsql

import "gorm.io/gorm"

type Clause interface {
	Clause(tx *gorm.DB) (*gorm.DB, error)
}
