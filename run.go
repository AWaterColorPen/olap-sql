package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

// RunChan executes the prepared *gorm.DB query and returns results over a channel.
// Each row is scanned into a map[string]any and sent on the returned channel.
// The channel is closed automatically when all rows have been consumed.
// This is suitable for streaming large result sets without loading them all into memory.
func RunChan(db *gorm.DB) (chan map[string]any, error) {
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}

	ch := make(chan map[string]any)
	go func() {
		defer close(ch)
		defer rows.Close()
		cnt := 0
		for rows.Next() {
			cnt++
			result := map[string]any{}
			_ = db.ScanRows(rows, &result)
			ch <- result
		}
	}()
	return ch, nil
}

// RunSync executes the prepared *gorm.DB query and returns all rows as a slice.
// Prefer RunSync for small-to-medium result sets; use RunChan for very large ones.
func RunSync(db *gorm.DB) ([]map[string]any, error) {
	var result []map[string]any
	return result, db.Scan(&result).Error
}

// BuildResultChan collects rows from a streaming channel and assembles a Result.
// It sets the dimension/metric names from the query and appends each incoming row.
func BuildResultChan(query *types.Query, in chan map[string]any) (*types.Result, error) {
	result := &types.Result{}
	result.SetDimensions(query)
	for v := range in {
		if err := result.AddSource(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// BuildResultSync wraps a slice of rows into a Result with dimension/metric metadata.
func BuildResultSync(query *types.Query, in []map[string]any) (*types.Result, error) {
	result := &types.Result{}
	result.SetDimensions(query)
	result.Source = in
	return result, nil
}
