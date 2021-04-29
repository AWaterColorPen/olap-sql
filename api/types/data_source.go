package types

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

type DataSourceType string

func (d DataSourceType) ToEnum() proto.DATA_SOURCE_TYPE {
	if v, ok := proto.DATA_SOURCE_TYPE_value[string(d)]; ok {
		return proto.DATA_SOURCE_TYPE(v)
	}
	return proto.DATA_SOURCE_TYPE_DATA_SOURCE_UNKNOWN
}

func EnumToDataSourceType(d proto.DATA_SOURCE_TYPE) DataSourceType {
	return DataSourceType(d.String())
}

const (
	DataSourceTypeUnknown    DataSourceType = "DATA_SOURCE_UNKNOWN"
	DataSourceTypeClickHouse DataSourceType = "DATA_SOURCE_CLICKHOUSE"
	DataSourceTypeDruid      DataSourceType = "DATA_SOURCE_DRUID"
	DataSourceTypeKylin      DataSourceType = "DATA_SOURCE_KYLIN"
	DataSourceTypePresto     DataSourceType = "DATA_SOURCE_PRESTO"
)

type DataSource struct {
	Type       DataSourceType `json:"type"`
	Name       string         `json:"name"`
	SubRequest *Request       `json:"sub_request"`
}

func (d *DataSource) Statement() (string, error) {
	if d.SubRequest == nil {
		return fmt.Sprintf("%v", d.Name), nil
	}

	statement, err := d.SubRequest.Statement()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("( %v ) %v", statement, d.Name), nil
}
