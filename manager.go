package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/awatercolorpen/olap-sql/dictionary"
	"gorm.io/gorm"
)

type Manager struct {
	clients    Clients
	dictionary *dictionary.Dictionary
}

func (m *Manager) GetClients() (Clients, error) {
	if m.clients == nil {
		return nil, fmt.Errorf("it is no initialization")
	}
	return m.clients, nil
}

func (m *Manager) GetDictionary() (*dictionary.Dictionary, error) {
	if m.dictionary == nil {
		return nil, fmt.Errorf("it is no initialization")
	}
	return m.dictionary, nil
}

func (m *Manager) RunSync(query *types.Query) (*types.Result, error) {
	db, err := m.BuildTransaction(query)
	if err != nil {
		return nil, err
	}
	rows, err := RunSync(db)
	if err != nil {
		return nil, err
	}
	return BuildResultSync(query, rows)
}

func (m *Manager) RunChan(query *types.Query) (*types.Result, error) {
	db, err := m.BuildTransaction(query)
	if err != nil {
		return nil, err
	}
	rows, err := RunChan(db)
	if err != nil {
		return nil, err
	}
	return BuildResultChan(query, rows)
}

func (m *Manager) BuildTransaction(query *types.Query) (*gorm.DB, error) {
	query.TranslateTimeIntervalToFilter()
	dict, err := m.GetDictionary()
	if err != nil {
		return nil, err
	}
	request, err := dict.Translate(query)
	if err != nil {
		return nil, err
	}
	clients, err := m.GetClients()
	if err != nil {
		return nil, err
	}
	return clients.SubmitClause(request)
}

func NewManager(configuration *Configuration) (*Manager, error) {
	// set default olap-sql option
	if configuration.DefaultParallelNumber != 0 {
		DefaultParallelNumber = configuration.DefaultParallelNumber
	}

	m := &Manager{}
	if configuration.ClientsOption != nil {
		clients, err := NewClients(configuration.ClientsOption)
		if err != nil {
			return nil, err
		}
		m.clients = clients
	}
	if configuration.DictionaryOption != nil {
		dict, err := dictionary.NewDictionary(configuration.DictionaryOption)
		if err != nil {
			return nil, err
		}
		m.dictionary = dict
	}
	return m, nil
}
