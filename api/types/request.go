package types

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Request struct {
	DBType     DBType       `json:"db_type"`
	Dataset    string       `json:"dataset"`
	Metrics    []*Metric    `json:"metrics"`
	Dimensions []*Dimension `json:"dimensions"`
	Filters    []*Filter    `json:"filters"`
	Joins      []*Join      `json:"joins"`
	Orders     []*OrderBy   `json:"orders"`
	Limit      *Limit       `json:"limit"`
	DataSource *DataSource  `json:"data_source"`
	Sql        string       `json:"sql"`
}

func (r *Request) GetDBType() DBType {
	return r.DBType
}

func (r *Request) GetDataset() string {
	return r.Dataset
}

func (r *Request) BuildDB(tx *gorm.DB) (*gorm.DB, error) {
	if err := r.SourceStatement(tx); err != nil {
		return nil, err
	}

	select1, err := r.dimensionStatement()
	if err != nil {
		return nil, err
	}
	select2, err := r.metricStatement()
	if err != nil {
		return nil, err
	}
	select3 := append([]string{}, select1...)
	select3 = append(select3, select2...)
	tx = tx.Select(select3)

	where1, err := r.filterStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range where1 {
		tx = tx.Where(v)
	}

	join1, err := r.joinStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range join1 {
		tx = tx.Joins(v)
	}

	group1, err := r.groupStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range group1 {
		tx = tx.Group(v)
	}
	order1, err := r.orderStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range order1 {
		tx = tx.Order(v)
	}

	table1, err := r.tableStatement()
	if err != nil {
		return nil, err
	}
	tx = tx.Table(table1)

	if r.Limit != nil {
		if r.Limit.Limit != 0 {
			tx = tx.Limit(int(r.Limit.Limit))
		}
		if r.Limit.Offset != 0 {
			tx = tx.Offset(int(r.Limit.Offset))
		}
	}

	return tx, nil
}

func (r *Request) BuildSQL(tx *gorm.DB) (string, error) {
	db, err := r.BuildDB(tx.Session(&gorm.Session{DryRun: true}))
	if err != nil {
		return "", err
	}
	_ = db.Scan(nil)
	stmt := db.Statement
	return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...), nil
}

func (r *Request) metricStatement() ([]string, error) {
	var statement []string
	for _, v := range r.Metrics {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (r *Request) dimensionStatement() ([]string, error) {
	var statement []string
	for _, v := range r.Dimensions {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (r *Request) filterStatement() ([]string, error) {
	var statement []string
	for _, v := range r.Filters {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (r *Request) joinStatement() ([]string, error) {
	if r.DataSource == nil {
		return nil, fmt.Errorf("nil data source")
	}

	var statement []string
	for _, v := range r.Joins {
		var on []string
		onName1, _ := v.DataSource1.GetDataSourceForOn()
		onName2, _ := v.DataSource2.GetDataSourceForOn()
		for _, u := range v.On {
			on = append(on, fmt.Sprintf("`%v`.`%v` = `%v`.`%v`", onName1, u.Key1, onName2, u.Key2))
		}

		joinName2, _ := v.DataSource2.GetDataSourceForJoin()

		switch r.DBType {
		case DBTypeSQLite, DBTypeClickHouse:
			statement = append(statement, fmt.Sprintf("LEFT JOIN %v ON %v", joinName2, strings.Join(on, " AND ")))
		default:
			return nil, fmt.Errorf("not supported db type %v", r.DBType)
		}
	}
	return statement, nil
}

func (r *Request) groupStatement() ([]string, error) {
	var statement []string
	for _, v := range r.Dimensions {
		s, err := v.Expression()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (r *Request) orderStatement() ([]string, error) {
	var statement []string
	for _, v := range r.Orders {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (r *Request) tableStatement() (string, error) {
	return r.DataSource.GetDataSourceForJoin()
}

func (r *Request) SourceStatement(tx *gorm.DB) error {
	if err := r.DataSource.Statement(tx); err != nil {
		return err
	}
	for _, join := range r.Joins {
		if err := join.DataSource1.Statement(tx); err != nil {
			return err
		}
		if err := join.DataSource2.Statement(tx); err != nil {
			return err
		}
	}
	return nil
}
