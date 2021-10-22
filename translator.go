package olapsql

import "github.com/awatercolorpen/olap-sql/api/types"

type Translator interface {
	Translate(*types.Query) (Clause, error)
}
