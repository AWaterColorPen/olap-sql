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
	Database  string         `json:"database"`
	Name      string         `json:"name"`
	AliasName string         `json:"alias"`
	Type      DataSourceType `json:"type"`
	// Joins    []*Join        `json:"joins"`
	Clause Clause `json:"clause"`

	expression string
}

func (d *DataSource) Expression() (string, error) {
	if d.expression != "" {
		return d.expression, nil
	}
	if d.Database == "" {
		return fmt.Sprintf("`%v`", d.Name), nil
	}
	return fmt.Sprintf("`%v`.`%v`", d.Database, d.Name), nil
}

func (d *DataSource) Alias() (string, error) {
	if d.AliasName == "" {
		return d.Name, nil
	}
	return d.AliasName, nil
}

func (d *DataSource) Statement() (string, error) {
	expression, _ := d.Expression()
	alias, _ := d.Alias()
	return fmt.Sprintf("(%v) AS %v", expression, alias), nil
}

func (d *DataSource) Init(tx *gorm.DB) (err error) {
	switch d.Clause {
	case nil:
		d.expression, err = d.Expression()
	default:
		d.expression, err = d.Clause.BuildSQL(tx)
	}
	return
}
