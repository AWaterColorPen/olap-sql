package olapsql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/awatercolorpen/olap-sql/api/types"
)

type parsing struct {
	columns []string
	values  []interface{}
}

func ParseChan(rows *sql.Rows) interface{} {
	in := make(chan interface{})
	go func() {
		defer close(in)
		defer rows.Close()
		// t1 := time.Now()

		ct, _ := rows.ColumnTypes()
		typeList := make([]reflect.Type, len(ct))
		for i, v := range ct {
			typeList[i] = v.ScanType()
		}

		byteType := reflect.TypeOf([]byte(""))

		cnt := 0
		for rows.Next() {
			cnt++
			cols, _ := rows.Columns()
			pointers := make([]interface{}, len(cols))
			for i := range pointers {
				if ct[i].DatabaseTypeName() == "double" {
					pointers[i] = reflect.New(typeList[i]).Interface()
				} else {
					pointers[i] = reflect.New(byteType).Interface()
				}
			}
			_ = rows.Scan(pointers...)
			in <- &parsing{columns: cols, values: pointers}
		}

		// latency := time.Since(t1).Microseconds()
		// log.Info("result len: %v, latency: %v ms", cnt, latency)
	}()

	ch := Parallel(in, func(v interface{}) interface{} {
		u := v.(*parsing)
		return ParseOneRow(u.columns, u.values)
	}, DefaultParallelNumber)
	return ch
}

func ParseSync(rows *sql.Rows) interface{} {
	// TODO
	return nil
}

func ParseOneRow(columns []string, values []interface{}) types.Item {
	item := map[string]string{}
	for i := range values {
		k := columns[i]
		switch w := values[i].(type) {
		case *[]byte:
			item[k] = string(*w)
		default:
			item[k] = fmt.Sprint(w)
		}
	}
	return item
}

