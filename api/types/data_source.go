package types

import (
	"fmt"
	"gorm.io/gorm"
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
	Database string         `json:"database"`
	Name     string         `json:"name"`
	Alias    string         `json:"alias"`
	Joins    []*Join        `json:"joins"`
	Request  *Request       `json:"sub_request"`
	Sql 	 string			`json:"sql"`
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

func (d *DataSource) Statement(tx *gorm.DB) error {
	switch d.Request {
	case nil:
		sql, err := d.getName()
		d.Sql = sql
		return err
	default:
		sql, err := d.Request.BuildSQL(tx)
		d.Sql = sql
		return err
	}
}

func (d *DataSource) GetDataSourceForOn() (string, error) {
	if d.Alias != "" {
		return d.Alias, nil
	}
	return d.Name, nil
}

func (d *DataSource) GetDataSourceForJoin() (string, error) {
	if d.Alias != "" {
		return fmt.Sprintf("(%v) AS %v", d.Sql, d.Alias), nil
	}
	return d.getName()
}