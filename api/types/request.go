package types

import (
	"fmt"
	"strings"
)

type Request struct {
	Metrics    []*Metric    `json:"metrics"`
	Dimensions []*Dimension `json:"dimensions"`
	Filters    []*Filter    `json:"filters"`
	Joins      []*Join      `json:"jsons"`
	DataSource *DataSource  `json:"data_source"`
}

func (r *Request) Statement() (string, error) {
	statement1, err := r.metricStatement()
	if err != nil {
		return "", err
	}

	statement2, err := r.dimensionStatement()
	if err != nil {
		return "", err
	}

	statement3, err := r.filterStatement()
	if err != nil {
		return "", err
	}

	statement4, err := r.joinStatement()
	if err != nil {
		return "", err
	}

	statement5, err := r.tableStatement()
	if err != nil {
		return "", err
	}

	return r.buildSql(statement1, statement2, statement3, statement4, statement5), nil
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
		return nil, ErrInvalidDataSource
	}

	var statement []string
	for _, v := range r.Joins {
		switch r.DataSource.Type {
		case DataSourceTypeKylin, DataSourceTypePresto, DataSourceTypeClickHouse:
			statement = append(statement, fmt.Sprintf("LEFT JOIN %v ON %v.%v = %v.%v", v.Table2, v.Table1, v.Key1, v.Table2, v.Key2))
		// case DataSourceTypeClickHouse:
		// 	statement = append(statement, fmt.Sprintf("t1 LEFT JOIN %v ON t1.%v = %v.%v", v.Table2, v.Key1, v.Table2, v.Key2))
		default:
			return nil, ErrNotSupportedDataSourceType
		}
	}
	return statement, nil
}

func (r *Request) tableStatement() (string, error) {
	return r.DataSource.Statement()
}

func (r *Request) buildSql(metrics, dimensions, filters, joins []string, table string) string {
	selectCol := append([]string{}, metrics...)
	selectCol = append(selectCol, dimensions...)

	selectStatement := strings.Join(selectCol, " , ")
	groupStatement := strings.Join(dimensions, " , ")
	whereStatement := strings.Join(filters, " AND ")
	joinStatement := strings.Join(joins, " ")

	sql := fmt.Sprintf("SELECT %v FROM %v", selectStatement, table)
	if joinStatement != "" {
		sql = fmt.Sprintf("%v %v", sql, joinStatement)
	}
	if whereStatement != "" {
		sql = fmt.Sprintf("%v WHERE %v", sql, whereStatement)
	}
	if groupStatement != "" {
		sql = fmt.Sprintf("%v GROUP BY %v", sql, groupStatement)
	}

	return sql
}
