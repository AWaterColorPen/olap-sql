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
			result := map[string]interface{} {}
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

func BuildTableResultChan(query *types.Query, in chan map[string]interface{}) (*types.TableResult, error) {
	result := &types.TableResult{}
	result.SetHeader(query)
	for v := range in {
		if err := result.AddRow(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func BuildTableResultSync(query *types.Query, in []map[string]interface{}) (*types.TableResult, error) {
	result := &types.TableResult{}
	result.SetHeader(query)
	for _, v := range in {
		if err := result.AddRow(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}
