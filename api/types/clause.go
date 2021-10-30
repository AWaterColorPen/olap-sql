package types

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Clause interface {
	GetDBType() DBType
	GetDataset() string
	BuildDB(tx *gorm.DB) (*gorm.DB, error)
	BuildSQL(tx *gorm.DB) (string, error)
}

type baseClause struct {
	DBType  DBType
	Dataset string
}

func (b *baseClause) GetDBType() DBType {
	return b.DBType
}

func (b *baseClause) GetDataset() string {
	return b.Dataset
}

type SqlClause struct {
	baseClause
	Sql string
}

func (s *SqlClause) BuildDB(tx *gorm.DB) (*gorm.DB, error) {
	return tx.Raw(s.Sql), nil
}

func (s *SqlClause) BuildSQL(*gorm.DB) (string, error) {
	return s.Sql, nil
}

type NormalClause struct {
	baseClause
	Metrics    []*Metric    `json:"metrics"`
	Dimensions []*Dimension `json:"dimensions"`
	Filters    []*Filter    `json:"filters"`
	Joins      []*Join      `json:"joins"`
	Orders     []*OrderBy   `json:"orders"`
	Limit      *Limit       `json:"limit"`
	DataSource *DataSource  `json:"data_source"`
}

func (n *NormalClause) BuildDB(tx *gorm.DB) (*gorm.DB, error) {
	if err := n.SourceStatement(tx); err != nil {
		return nil, err
	}

	select1, err := n.dimensionStatement()
	if err != nil {
		return nil, err
	}
	select2, err := n.metricStatement()
	if err != nil {
		return nil, err
	}
	select3 := append([]string{}, select1...)
	select3 = append(select3, select2...)
	tx = tx.Select(select3)

	where1, err := n.filterStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range where1 {
		tx = tx.Where(v)
	}

	join1, err := n.joinStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range join1 {
		tx = tx.Joins(v)
	}

	group1, err := n.groupStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range group1 {
		tx = tx.Group(v)
	}
	order1, err := n.orderStatement()
	if err != nil {
		return nil, err
	}
	for _, v := range order1 {
		tx = tx.Order(v)
	}

	table1, err := n.tableStatement()
	if err != nil {
		return nil, err
	}
	tx = tx.Table(table1)

	if n.Limit != nil {
		if n.Limit.Limit != 0 {
			tx = tx.Limit(int(n.Limit.Limit))
		}
		if n.Limit.Offset != 0 {
			tx = tx.Offset(int(n.Limit.Offset))
		}
	}

	return tx, nil
}

func (n *NormalClause) BuildSQL(tx *gorm.DB) (string, error) {
	db, err := n.BuildDB(tx.Session(&gorm.Session{DryRun: true}))
	if err != nil {
		return "", err
	}
	_ = db.Scan(nil)
	stmt := db.Statement
	return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...), nil
}

func (n *NormalClause) metricStatement() ([]string, error) {
	var statement []string
	for _, v := range n.Metrics {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (n *NormalClause) dimensionStatement() ([]string, error) {
	var statement []string
	for _, v := range n.Dimensions {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (n *NormalClause) filterStatement() ([]string, error) {
	var statement []string
	for _, v := range n.Filters {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (n *NormalClause) joinStatement() ([]string, error) {
	if n.DataSource == nil {
		return nil, fmt.Errorf("nil data source")
	}

	var statement []string
	for _, v := range n.Joins {
		var on []string
		onName1, _ := v.DataSource1.GetDataSourceForOn()
		onName2, _ := v.DataSource2.GetDataSourceForOn()
		for _, u := range v.On {
			on = append(on, fmt.Sprintf("`%v`.`%v` = `%v`.`%v`", onName1, u.Key1, onName2, u.Key2))
		}

		joinName2, _ := v.DataSource2.GetDataSourceForJoin()

		switch n.DBType {
		case DBTypeSQLite, DBTypeClickHouse:
			statement = append(statement, fmt.Sprintf("LEFT JOIN %v ON %v", joinName2, strings.Join(on, " AND ")))
		default:
			return nil, fmt.Errorf("not supported db type %v", n.DBType)
		}
	}
	return statement, nil
}

func (n *NormalClause) groupStatement() ([]string, error) {
	var statement []string
	for _, v := range n.Dimensions {
		s, err := v.Expression()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (n *NormalClause) orderStatement() ([]string, error) {
	var statement []string
	for _, v := range n.Orders {
		s, err := v.Statement()
		if err != nil {
			return nil, err
		}
		statement = append(statement, s)
	}
	return statement, nil
}

func (n *NormalClause) tableStatement() (string, error) {
	return n.DataSource.GetDataSourceForJoin()
}

func (n *NormalClause) SourceStatement(tx *gorm.DB) error {
	if err := n.DataSource.Statement(tx); err != nil {
		return err
	}
	for _, join := range n.Joins {
		if err := join.DataSource1.Statement(tx); err != nil {
			return err
		}
		if err := join.DataSource2.Statement(tx); err != nil {
			return err
		}
	}
	return nil
}
