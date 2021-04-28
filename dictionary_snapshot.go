package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/models"
)

type DictionarySnapshot struct {
	metrics     []*models.Metric
	dimensions  []*models.Dimension
	dataSources []*models.DataSource
	dataSets    []*models.DataSet
}

// func (d *DictionarySnapshot) Translate(in interface{}) (interface{}, error) {
// 	query, ok := in.(*types.Query)
// 	if !ok {
// 		return nil, nil
// 	}
//
// 	query.DataSet
//
// 	request := &types.Request{}
// 	return request, nil
// }