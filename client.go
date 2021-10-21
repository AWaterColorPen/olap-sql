package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
)

type ClientsOption = map[string]*DBOption

type Clients map[string]*gorm.DB

func (c Clients) RegisterByKV(dbType types.DBType, dataset string, db *gorm.DB) {
	key := c.key(dbType, dataset)
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

func (c Clients) Get(dbType types.DBType, dataset string) (*gorm.DB, error) {
	key1 := c.key(dbType, dataset)
	if v, ok := c[key1]; ok {
		return v, nil
	}
	key2 := c.key(dbType, "")
	if v, ok := c[key2]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found client %v %v", dbType, dataset)
}

func (c Clients) key(dbType types.DBType, dataset string) string {
	if dataset == "" {
		return fmt.Sprintf("%v", dbType)
	}
	return fmt.Sprintf("%v/%v", dbType, dataset)
}

func (c Clients) BuildDB(clause Clause) (*gorm.DB, error) {
	client, err := c.Get(clause.GetDBType(), clause.GetDataset())
	if err != nil {
		return nil, err
	}
	return clause.BuildDB(client)
}

func (c Clients) BuildSQL(clause Clause) (string, error) {
	client, err := c.Get(clause.GetDBType(), clause.GetDataset())
	if err != nil {
		return "", err
	}
	return clause.BuildSQL(client)
}

func NewClients(option ClientsOption) (Clients, error) {
	c := Clients{}
	if err := c.RegisterByOption(option); err != nil {
		return nil, err
	}
	return c, nil
}
