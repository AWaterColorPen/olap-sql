package types

import (
	"fmt"
)

type DataSourceType string

const (
	DataSourceTypeUnknown    DataSourceType = "DATA_SOURCE_UNKNOWN"
	DataSourceTypeClickHouse DataSourceType = "DATA_SOURCE_CLICKHOUSE"
	DataSourceTypeDruid      DataSourceType = "DATA_SOURCE_DRUID"
	DataSourceTypeKylin      DataSourceType = "DATA_SOURCE_KYLIN"
	DataSourceTypePresto     DataSourceType = "DATA_SOURCE_PRESTO"
)

type DataSource struct {
	Type DataSourceType `json:"type"`
	Name string         `json:"name"`
	// SubRequest *Request       `json:"sub_request"`
}

func (d *DataSource) Statement() (string, error) {
	return fmt.Sprintf("%v", d.Name), nil
	// if d.SubRequest == nil {
	// 	return fmt.Sprintf("%v", d.Name), nil
	// }

	// statement, err := d.SubRequest.Statement()
	// if err != nil {
	// 	return "", err
	// }
	// return fmt.Sprintf("( %v ) %v", statement, d.Name), nil
}
