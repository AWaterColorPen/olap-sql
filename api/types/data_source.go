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
	Type     DataSourceType `json:"type"`
	Database string         `json:"database"`
	Name     string         `json:"name"`
	Alias    string         `json:"alias"`
	Joins    []*Join        `json:"joins"`
	Request  *Request       `json:"sub_request"`
}

func (d *DataSource) getName() (string, error) {
	if d.Database == "" {
		return fmt.Sprintf("`%v`", d.Name), nil
	}
	return fmt.Sprintf("`%v`.`%v`", d.Database, d.Name), nil
}

func (d *DataSource) getAlias() (string, error) {
	if d.Alias == "" {
		return "", fmt.Errorf("alias is nil")
	}
	return d.Alias, nil
}

func (d *DataSource) Statement() (string, error) {
	switch d.Request {
	// 寻常的datasource
	case nil:
		return d.getName()
	// select子查询
	default:
		return d.Request.BuildSQL()
	}
}

func (d *DataSource) GetDataSourceForOn() (string, error) {
	// 对于on来说 要么用别名(也就是select的情况)，不然的话就是普通的情况
	if d.Alias != "" {
		return d.Alias, nil
	}
	return d.Name, nil
}

func (d *DataSource) GetDataSourceForJoin() (string, error) {
	return "", nil
}