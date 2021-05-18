package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

func RunChan(db *gorm.DB) (chan map[string]interface{}, error) {
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}

	ch := make(chan map[string]interface{})
	go func() {
		defer close(ch)
		defer rows.Close()
		cnt := 0
		for rows.Next() {
			cnt++
			result := map[string]interface{}{}
			_ = db.ScanRows(rows, &result)
			ch <- result
		}
	}()
	return ch, nil
}

func RunSync(db *gorm.DB) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	return result, db.Scan(&result).Error
}

func BuildResultChan(query *types.Query, in chan map[string]interface{}) (*types.Result, error) {
	result := &types.Result{}
	result.SetDimensions(query)
	for v := range in {
		if err := result.AddSource(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func BuildResultSync(query *types.Query, in []map[string]interface{}) (*types.Result, error) {
	result := &types.Result{}
	result.SetDimensions(query)
	for _, v := range in {
		if err := result.AddSource(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}
