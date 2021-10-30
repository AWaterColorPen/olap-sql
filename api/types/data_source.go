package types

import (
	"fmt"
	"gorm.io/gorm"
)

type DataSourceType string

const (
	DataSourceTypeFact              DataSourceType = "fact"
	DataSourceTypeDimension         DataSourceType = "dimension"
	DataSourceTypeFactDimensionJoin DataSourceType = "fact_dimension_join"
	DataSourceTypeMergedJoin        DataSourceType = "merged_join"
)

type DataSource struct {
	Database string         `json:"database"`
	Name     string         `json:"name"`
	Alias    string         `json:"alias"`
	Type     DataSourceType `json:"type"`
	Joins    []*Join        `json:"joins"`
	Clause   Clause         `json:"clause"`
	Sql      string         `json:"Sql"`
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
	switch d.Clause {
	case nil:
		sql, err := d.getName()
		d.Sql = sql
		return err
	default:
		sql, err := d.Clause.BuildSQL(tx)
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
