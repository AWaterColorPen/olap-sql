package olapsql

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type ClientsOption = map[string]*DBOption

type Clients map[string]*gorm.DB

func (c Clients) RegisterByKV(dataSourceType types.DataSourceType, dataset string, db *gorm.DB) {
	key := c.key(dataSourceType, dataset)
	c[key] = db
}

func (c Clients) RegisterByOption(option ClientsOption) error {
	for k, v := range option {
		db, err := v.NewDB()
		if err != nil {
			return err
		}
		c[k] = db
	}
	return nil
}

func (c Clients) Get(dataSourceType types.DataSourceType, dataset string) (*gorm.DB, error) {
	key1 := c.key(dataSourceType, dataset)
	if v, ok := c[key1]; ok {
		return v, nil
	}
	key2 := c.key(dataSourceType, "")
	if v, ok := c[key2]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found client %v %v", dataSourceType, dataset)
}

func (c Clients) key(dataSourceType types.DataSourceType, dataset string) string {
	if dataset == "" {
		return fmt.Sprintf("%v", dataSourceType)
	}
	return fmt.Sprintf("%v/%v", dataSourceType, dataset)
}

func (c Clients) SubmitRequest(request *types.Request) (*types.Response, error) {
	client, err := c.Get(request.DataSource.Type, request.DataSource.Name)
	if err != nil {
		return nil, err
	}
	db, err := request.Clause(client)
	if err != nil {
		return nil, err
	}
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	ch := ParseChan(rows)
	response := &types.Response{}
	linq.From(ch).ForEach(func(v interface{}) {
		u := v.(*types.Item)
		response.Rows = append(response.Rows, u)
	})
	return response, nil
}

func (c Clients) SubmitSql(request *types.Request) (*types.Response, error) {
	client, err := c.Get(request.DataSource.Type, request.DataSource.Name)
	if err != nil {
		return nil, err
	}
	db, err := request.Clause(client)
	if err != nil {
		return nil, err
	}
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	ch := ParseChan(rows)
	response := &types.Response{}
	linq.From(ch).ForEach(func(v interface{}) {
		u := v.(*types.Item)
		response.Rows = append(response.Rows, u)
	})
	return response, nil
}

func NewClients(option ClientsOption) (Clients, error) {
	c := Clients{}
	if err := c.RegisterByOption(option); err != nil {
		return nil, err
	}
	return c, nil
}